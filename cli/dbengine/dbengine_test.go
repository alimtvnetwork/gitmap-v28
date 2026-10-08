package dbengine

import (
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

type TestItem struct {
	ItemId   uint64
	ItemName string
	Category string
	IsActive bool
}

type TestItemFieldType string

func (e TestItemFieldType) Name() string   { return string(e) }
func (e TestItemFieldType) String() string { return string(e) }
func (e TestItemFieldType) Value() string  { return string(e) }
func (e TestItemFieldType) IsCompare(target TestItemFieldType) bool {
	return e == target
}

func (e TestItemFieldType) IsEnum() bool {
	return testItemValidMap[e]
}

func (e TestItemFieldType) IsItemId() bool   { return e == TestItemDb.ItemId }
func (e TestItemFieldType) IsItemName() bool { return e == TestItemDb.ItemName }
func (e TestItemFieldType) IsCategory() bool { return e == TestItemDb.Category }
func (e TestItemFieldType) IsIsActive() bool { return e == TestItemDb.IsActive }

func (e TestItemFieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

func (e *TestItemFieldType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	target := TestItemFieldType(s)
	if !testItemValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid test item enum: %s", s), "unmarshal test item field")
	}

	*e = target

	return nil
}

func (e TestItemFieldType) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(string(e))
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize field to json")
	}

	return string(b), nil
}

func (e *TestItemFieldType) FromJSON(s string) *apperror.AppError {
	var str string
	if err := json.Unmarshal([]byte(s), &str); err != nil {
		return apperror.WrapSimple(err, "deserialize field from json")
	}

	target := TestItemFieldType(str)
	if !testItemValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid test item enum: %s", str), "validate field from json")
	}

	*e = target

	return nil
}

type testItemDbRegistry struct {
	ItemId   TestItemFieldType
	ItemName TestItemFieldType
	Category TestItemFieldType
	IsActive TestItemFieldType
}

func (r testItemDbRegistry) All() []TestItemFieldType {
	return []TestItemFieldType{r.ItemId, r.ItemName, r.Category, r.IsActive}
}

func (r testItemDbRegistry) Names() []string {
	return []string{"ItemId", "ItemName", "Category", "IsActive"}
}

func (r testItemDbRegistry) IsEnum(target TestItemFieldType) bool {
	return testItemValidMap[target]
}

func (r testItemDbRegistry) IsItemId(target TestItemFieldType) bool {
	return target == r.ItemId
}

func (r testItemDbRegistry) IsItemName(target TestItemFieldType) bool {
	return target == r.ItemName
}

func (r testItemDbRegistry) IsCategory(target TestItemFieldType) bool {
	return target == r.Category
}

func (r testItemDbRegistry) IsIsActive(target TestItemFieldType) bool {
	return target == r.IsActive
}

func (r testItemDbRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize test item db registry to json")
	}

	return string(b), nil
}

var TestItemDb = testItemDbRegistry{
	ItemId:   "ItemId",
	ItemName: "ItemName",
	Category: "Category",
	IsActive: "IsActive",
}

var testItemValidMap = map[TestItemFieldType]bool{
	TestItemDb.ItemId:   true,
	TestItemDb.ItemName: true,
	TestItemDb.Category: true,
	TestItemDb.IsActive: true,
}

var TestItemField = TestItemDb

func scanTestItem(s RowScanner) (*TestItem, error) {
	var item TestItem
	var activeInt int
	err := s.Scan(&item.ItemId, &item.ItemName, &item.Category, &activeInt)
	if err != nil {
		return nil, err
	}

	item.IsActive = activeInt == 1

	return &item, nil
}
