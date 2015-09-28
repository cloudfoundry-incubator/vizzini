package vizzini_test

import (
	"github.com/cloudfoundry-incubator/bbs/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cells", func() {
	It("should return all cells", func() {
		cells, err := serviceClient.Cells(logger)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(len(cells)).Should(BeNumerically(">=", 1))

		var cell_z1_0 *models.CellPresence
		for _, cell := range cells {
			if cell.CellID == "cell_z1-0" {
				cell_z1_0 = cell
				break
			}
		}

		Ω(cell_z1_0).ShouldNot(BeNil())
		Ω(cell_z1_0.CellID).Should(Equal("cell_z1-0"))
		Ω(cell_z1_0.Zone).Should(Equal("z1"))
		Ω(cell_z1_0.Capacity.MemoryMB).Should(BeNumerically(">", 0))
		Ω(cell_z1_0.Capacity.DiskMB).Should(BeNumerically(">", 0))
		Ω(cell_z1_0.Capacity.Containers).Should(BeNumerically(">", 0))
	})
})
