package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xtreme-andleung/whiteboardbot/app"
	"github.com/nlopes/slack"
	"github.com/xtreme-andleung/whiteboardbot/spec"
)

var _ = Describe("Faces Integration", func() {
	var (
		newFaceEvent slack.MessageEvent
		setNameEvent slack.MessageEvent
		setDateEvent slack.MessageEvent
		client spec.MockSlackClient
		clock spec.MockClock
	)

	BeforeEach(func() {
		client = spec.MockSlackClient{}
		clock = spec.MockClock{}

		newFaceEvent = slack.MessageEvent{}
		newFaceEvent.Text = "wb faces"

		setNameEvent = slack.MessageEvent{}
		setNameEvent.Text = "wb name Dariusz Lorenc"

		setDateEvent = slack.MessageEvent{}
		setDateEvent.Text = "wb date 2015-12-01"
	})

	Describe("with faces keyword", func() {
		It("should begin creating a new face entry and respond with face string", func() {
			_, Text := ParseMessageEvent(&client, clock, &newFaceEvent)
			Expect(Text).To(Equal("faces\n  *name: \n  date: 2015-01-02"))
		})
	})

	Context("setting a name detail", func() {
		Describe("with a new face entry started", func() {
			BeforeEach(func() {
				ParseMessageEvent(&client, clock, &newFaceEvent)
			})
			Describe("with correct keyword", func() {
				It("should set the name of the entry and respond with face string", func() {
					_, Text := ParseMessageEvent(&client, clock, &setNameEvent)
					Expect(Text).To(Equal("faces\n  *name: Dariusz Lorenc\n  date: 2015-01-02"))
				})
			})
			Describe("with incorrect keyword", func() {
				It("should respond with default", func() {
					setNameEvent.Text = "wb nameSomethingWrong"
					_, Text := ParseMessageEvent(&client, clock, &setNameEvent)
					Expect(Text).To(Equal("aleung no you nameSomethingWrong"))
				})
			})
		})
	})
	Context("setting a date detail", func() {
		Describe("with a new face entry started", func() {
			BeforeEach(func() {
				ParseMessageEvent(&client, clock, &newFaceEvent)
			})
			Describe("with correct keyword", func() {
				It("should set the date of the entry and respond with face string", func() {
					_, Text := ParseMessageEvent(&client, clock, &setDateEvent)
					Expect(Text).To(Equal("faces\n  *name: \n  date: 2015-12-01"))
				})
				It("should not set invalid date and respond with help message", func() {
					setDateEvent.Text = "wb date 12/01/2015"
					_, Text := ParseMessageEvent(&client, clock, &setDateEvent)
					Expect(Text).To(Equal("faces\n  *name: \n  date: 2015-01-02\nDate not set, use YYYY-MM-DD as date format"))
				})
			})
			Describe("with incorrect keyword", func() {
				It("should respond with default", func() {
					setDateEvent.Text = "wb date2015-12-01"
					_, Text := ParseMessageEvent(&client, clock, &setDateEvent)
					Expect(Text).To(Equal("aleung no you date2015-12-01"))
				})
			})
		})
	})
})


