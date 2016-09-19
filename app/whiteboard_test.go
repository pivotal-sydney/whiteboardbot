package app_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-sydney/whiteboardbot/app"
	"github.com/pivotal-sydney/whiteboardbot/model"
	"github.com/pivotal-sydney/whiteboardbot/spec"
)

var _ = Describe("Whiteboard", func() {

	var whiteboard app.WhiteboardApp

	BeforeEach(func() {
		slackClient := spec.MockSlackClient{}
		clock := spec.MockClock{}
		restClient := spec.MockRestClient{}
		store := spec.MockStore{}
		whiteboard = app.NewWhiteboard(&slackClient, &restClient, clock, &store)
	})

	Describe("Filter out old entries", func() {
		It("should return all entries within X days", func() {
			entries := []model.Entry{
				{Title: "today", Date: "2015-01-02", Author: "Andrew"},
				{Title: "in five days", Date: "2015-01-07", Author: "Andrew"},
				{Title: "in six days", Date: "2015-01-08", Author: "Dariusz"}}

			filteredEntries := whiteboard.FilterOutOld(entries, 5, "Australia/Sydney")

			Expect(filteredEntries).To(HaveLen(2))
			Expect(filteredEntries[0]).To(Equal(entries[0]))
			Expect(filteredEntries[1]).To(Equal(entries[1]))
		})
		It("should still return entries within invalid dates", func() {
			entries := []model.Entry{
				{Title: "empty date", Date: "", Author: "Andrew"},
				{Title: "invalid date", Date: "invalid date", Author: "Andrew"},
				{Title: "in six days", Date: "2015-01-08", Author: "Dariusz"}}

			filteredEntries := whiteboard.FilterOutOld(entries, 5, "Australia/Sydney")

			Expect(filteredEntries).To(HaveLen(2))
			Expect(filteredEntries[0]).To(Equal(entries[0]))
			Expect(filteredEntries[1]).To(Equal(entries[1]))
		})
	})
})
