package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const LOCAL_APP_URL = "http://localhost:8080"

func baseURL() string {
	if url := os.Getenv("API_URL"); url != "" {
		return url
	}
	return LOCAL_APP_URL
}

var _ = Describe("Optimizer API", func() {

	// in order to isolate each test, restore default sizes before each
	BeforeEach(func() {
		resp := putJSON("/packs/sizes", map[string]any{
			"sizes": []int{250, 500, 1000, 2000, 5000},
		})
		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
	})

	Describe("GET /packs/sizes", func() {
		It("returns the default pack sizes", func() {
			var result struct {
				Sizes []int `json:"sizes"`
			}
			resp := getJSON("/packs/sizes", &result)

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result.Sizes).To(ConsistOf(250, 500, 1000, 2000, 5000))
		})
	})

	Describe("PUT /packs/sizes", func() {
		It("updates pack sizes successfully", func() {
			resp := putJSON("/packs/sizes", map[string]any{
				"sizes": []int{100, 200, 300},
			})
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			var result struct {
				Sizes []int `json:"sizes"`
			}
			getJSON("/packs/sizes", &result)
			Expect(result.Sizes).To(ConsistOf(100, 200, 300))
		})

		It("returns 400 for an empty sizes array", func() {
			resp := putJSON("/packs/sizes", map[string]any{
				"sizes": []int{},
			})
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("returns 400 for non-positive size values", func() {
			resp := putJSON("/packs/sizes", map[string]any{
				"sizes": []int{-100, 500},
			})
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("returns 400 for zero size value", func() {
			resp := putJSON("/packs/sizes", map[string]any{
				"sizes": []int{0, 500},
			})
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("returns 400 when body is missing", func() {
			req, err := http.NewRequest(http.MethodPut, baseURL()+"/packs/sizes", nil)
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("GET /packs/calculate", func() {
		type packResult struct {
			Size     int `json:"size"`
			Quantity int `json:"quantity"`
		}
		type calcResponse struct {
			Packs []packResult `json:"packs"`
		}

		toMap := func(packs []packResult) map[int]int {
			m := make(map[int]int)
			for _, p := range packs {
				m[p.Size] = p.Quantity
			}
			return m
		}

		DescribeTable("calculates optimal packs",
			func(quantity int, expected map[int]int) {
				var result calcResponse
				resp := getJSON(fmt.Sprintf("/packs/calculate?quantity=%d", quantity), &result)

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(toMap(result.Packs)).To(Equal(expected))
			},
			Entry("1 item → 1×250", 1, map[int]int{250: 1}),
			Entry("250 items → 1×250", 250, map[int]int{250: 1}),
			Entry("251 items → 1×500", 251, map[int]int{500: 1}),
			Entry("500 items → 1×500", 500, map[int]int{500: 1}),
			Entry("501 items → 1×500 + 1×250", 501, map[int]int{500: 1, 250: 1}),
			Entry("12001 items → 2×5000 + 1×2000 + 1×250", 12001, map[int]int{5000: 2, 2000: 1, 250: 1}),
			Entry("10000 items → 2×5000", 10000, map[int]int{5000: 2}),
		)

		It("returns 400 when quantity is missing", func() {
			resp, err := http.Get(baseURL() + "/packs/calculate")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("returns 400 when quantity is zero", func() {
			resp, err := http.Get(baseURL() + "/packs/calculate?quantity=0")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("returns 400 when quantity is negative", func() {
			resp, err := http.Get(baseURL() + "/packs/calculate?quantity=-1")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("recalculates correctly after pack sizes are updated", func() {
			resp := putJSON("/packs/sizes", map[string]any{
				"sizes": []int{23, 31, 53},
			})
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			var result calcResponse
			getJSON(fmt.Sprintf("/packs/calculate?quantity=%d", 500000), &result)
			Expect(toMap(result.Packs)).To(Equal(map[int]int{23: 2, 31: 7, 53: 9429}))
		})
	})
})

func getJSON(path string, target any) *http.Response {
	resp, err := http.Get(baseURL() + path)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()
	Expect(json.NewDecoder(resp.Body).Decode(target)).To(Succeed())
	return resp
}

func putJSON(path string, body any) *http.Response {
	b, err := json.Marshal(body)
	Expect(err).NotTo(HaveOccurred())
	req, err := http.NewRequest(http.MethodPut, baseURL()+path, bytes.NewReader(b))
	Expect(err).NotTo(HaveOccurred())
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	Expect(err).NotTo(HaveOccurred())
	return resp
}
