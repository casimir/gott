package gott

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

var testData = []struct {
	in       string
	expected TaskData
}{
	// correct samples
	{in: "(A) Thank Mom for the meatballs @phone",
		expected: TaskData{
			Priority: 'A',
			Context:  []string{"phone"},
			Text:     "Thank Mom for the meatballs",
		}},
	{in: "(B) Schedule Goodwill pickup +GarageSale @phone",
		expected: TaskData{
			Priority: 'B',
			Project:  []string{"GarageSale"},
			Context:  []string{"phone"},
			Text:     "Schedule Goodwill pickup",
		},
	},
	{in: "Post signs around the neighborhood +GarageSale",
		expected: TaskData{
			Project: []string{"GarageSale"},
			Text:    "Post signs around the neighborhood",
		},
	},
	{in: "@GroceryStore Eskimo pies",
		expected: TaskData{
			Context: []string{"GroceryStore"},
			Text:    "Eskimo pies",
		},
	},
	{in: "2011-03-02 Document +TodoTxt task format",
		expected: TaskData{
			Date:    time.Date(2011, 3, 2, 0, 0, 0, 0, time.UTC),
			Project: []string{"TodoTxt"},
			Text:    "Document task format",
		}},
	{in: "(A) 2011-03-02 Call Mom",
		expected: TaskData{
			Priority: 'A',
			Date:     time.Date(2011, 3, 2, 0, 0, 0, 0, time.UTC),
			Text:     "Call Mom",
		}},
	{in: "(A) Call Mom +Family +PeaceLoveAndHappiness @iphone @phone",
		expected: TaskData{
			Priority: 'A',
			Project:  []string{"Family", "PeaceLoveAndHappiness"},
			Context:  []string{"iphone", "phone"},
			Text:     "Call Mom",
		}},
	{in: "x 2011-03-03 Call Mom",
		expected: TaskData{
			Done: true,
			Date: time.Date(2011, 3, 3, 0, 0, 0, 0, time.UTC),
			Text: "Call Mom",
		}},
	// invalid priority
	{in: "Really gotta call Mom (A) @phone @someday",
		expected: TaskData{
			Context: []string{"phone", "someday"},
			Text:    "Really gotta call Mom (A)",
		}},
	{in: "(B)->Submit TPS report",
		expected: TaskData{
			Text: "(B)->Submit TPS report",
		}},
	// invalid creation date
	{in: "(A) Call Mom 2011-03-02",
		expected: TaskData{
			Priority: 'A',
			Text:     "Call Mom 2011-03-02",
		}},
	// Not a project/context
	{in: "Email SoAndSo at soandso@example.com",
		expected: TaskData{
			Text: "Email SoAndSo at soandso@example.com",
		}},
	{in: "Learn how to add 2+2",
		expected: TaskData{
			Text: "Learn how to add 2+2",
		}},
	// Not done
	{in: "xylophone lesson",
		expected: TaskData{
			Text: "xylophone lesson",
		}},
	{in: "X 2012-01-01 Make resolutions",
		expected: TaskData{
			Text: "X 2012-01-01 Make resolutions",
		}},
	{in: "(A) x Find ticket prices",
		expected: TaskData{
			Priority: 'A',
			Text:     "x Find ticket prices",
		}},
}

func TestParsing(t *testing.T) {
	assert := assertions.New(t)
	for _, it := range testData {
		assert.So(NewTask(it.in), should.Resemble, it.expected)
	}
}

var niceList = TaskList{
	{Priority: 'A', Date: time.Date(2015, time.May, 2, 0, 0, 0, 0, time.UTC)},
	{Priority: 'A', Date: time.Date(2015, time.May, 12, 0, 0, 0, 0, time.UTC)},
	{Priority: 'A'},
	{Priority: 'B', Date: time.Date(2015, time.May, 2, 0, 0, 0, 0, time.UTC)},
	{Priority: 'B', Date: time.Date(2015, time.May, 2, 0, 0, 0, 0, time.UTC)},
	{Priority: 'C', Date: time.Date(2015, time.May, 1, 0, 0, 0, 0, time.UTC)},
	{Date: time.Date(2015, time.May, 1, 0, 0, 0, 0, time.UTC)},
	{Date: time.Date(2015, time.July, 1, 0, 0, 0, 0, time.UTC)},
	{},
	{},
	{},
}

func shuffle(ts TaskList) TaskList {
	rand.Seed(time.Now().UnixNano())
	shuffled := make(TaskList, len(ts))
	for i, r := range rand.Perm(len(ts)) {
		shuffled[i] = ts[r]
	}
	return shuffled
}

func TestSorting(t *testing.T) {
	assert := assertions.New(t)
	shuffled := shuffle(niceList)
	assert.So(shuffled, should.NotResemble, niceList)
	sort.Sort(shuffled)
	assert.So(shuffled, should.Resemble, niceList)
}
