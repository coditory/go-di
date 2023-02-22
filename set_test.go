package di

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SetSuite struct {
	suite.Suite
}

func (suite *SetSuite) TestEmptySetToSlice() {
	set := NewSet[string]()
	slice := set.ToSlice()
	suite.Equal([]string{}, slice)
}

func (suite *SetSuite) TestCreateSliceWithInsertionOrder() {
	set := NewSet[string]()
	set.Add("a")
	set.Add("b")
	set.Add("c")
	set.Add("a")
	slice := set.ToSlice()
	suite.Equal([]string{"b", "c", "a"}, slice)
}

func (suite *SetSuite) TestSetAutoGrow() {
	set := NewSet[string]()
	set.Add("a")
	set.Add("b")
	set.Add("c")
	slice := set.ToSlice()
	suite.Equal([]string{"a", "b", "c"}, slice)
}

func (suite *SetSuite) TestRemoveElement() {
	set := NewSet[string]()
	set.Add("a")
	set.Add("b")
	set.Add("c")
	set.Remove("b")
	suite.Equal([]string{"a", "c"}, set.ToSlice())
	suite.Equal(false, set.Contains("b"))
}

func TestSetSuite(t *testing.T) {
	suite.Run(t, new(SetSuite))
}
