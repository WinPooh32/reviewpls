package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/WinPooh32/opt"
)

type commonJSONSchema struct {
	Type        string  `json:"type,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Nullable    *bool   `json:"nullable,omitempty"`
}

type sortableProps struct {
	Names []string
	Props []any
}

func (x *sortableProps) Len() int {
	return len(x.Names)
}

func (x *sortableProps) Less(i, j int) bool {
	return x.Names[i] < x.Names[j]
}

func (x *sortableProps) Swap(i, j int) {
	x.Names[i], x.Names[j] = x.Names[j], x.Names[i]
	x.Props[i], x.Props[j] = x.Props[j], x.Props[i]
}

type propertyDef[T property, B any, PtrB *B] struct {
	prop T

	title       opt.T[string]
	description opt.T[string]
	nullable    opt.T[bool]
	enum        []PtrB
}

func newPropertyDef[T property, B any, PtrB *B](prop T) propertyDef[T, B, PtrB] {
	var def propertyDef[T, B, PtrB]

	def.prop = prop

	return def
}

func (p *propertyDef[T, B, PtrB]) Title(title string) T {
	p.title = opt.Wrap(title)
	return p.prop
}

func (p *propertyDef[T, B, PtrB]) Description(description string) T {
	p.description = opt.Wrap(description)
	return p.prop
}

func (p *propertyDef[T, B, PtrB]) Nullable(nullable bool) T {
	p.nullable = opt.Wrap(nullable)
	return p.prop
}

func (p *propertyDef[T, B, PtrB]) Enum(enum ...B) T {
	for _, v := range enum {
		p.enum = append(p.enum, &v)
	}

	return p.prop
}

func (p *propertyDef[T, B, PtrB]) EnumNull() T {
	var null PtrB

	p.enum = append(p.enum, null)

	return p.prop
}

type StringDefinition struct {
	propertyDef[*StringDefinition, string, *string]
}

func String() *StringDefinition {
	var def StringDefinition

	def.propertyDef = newPropertyDef[*StringDefinition, string](&def)

	return &def
}

func (sd *StringDefinition) MarshalJSON() ([]byte, error) {
	if sd == nil {
		return []byte("null"), nil
	}

	var value struct {
		commonJSONSchema
		Enum []*string `json:"enum,omitempty"`
	}

	value.Type = "string"

	if sd.title.Set() {
		value.Title = ptr(sd.title.Value())
	}

	if sd.description.Set() {
		value.Description = ptr(sd.description.Value())
	}

	if sd.nullable.Set() {
		value.Nullable = ptr(sd.nullable.Value())
	}

	value.Enum = sd.enum

	bs, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode string json: %w", err)
	}

	return bs, nil
}

type IntegerDefinition struct {
	propertyDef[*IntegerDefinition, int, *int]
}

func Integer() *IntegerDefinition {
	var def IntegerDefinition

	def.propertyDef = newPropertyDef[*IntegerDefinition, int](&def)

	return &def
}

func (id *IntegerDefinition) MarshalJSON() ([]byte, error) {
	if id == nil {
		return []byte("null"), nil
	}

	var value struct {
		commonJSONSchema
		Enum []*int `json:"enum,omitempty"`
	}

	value.Type = "integer"

	if id.title.Set() {
		value.Title = ptr(id.title.Value())
	}

	if id.description.Set() {
		value.Description = ptr(id.description.Value())
	}

	if id.nullable.Set() {
		value.Nullable = ptr(id.nullable.Value())
	}

	value.Enum = id.enum

	bs, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode integer json: %w", err)
	}

	return bs, nil
}

type NumberDefinition struct {
	propertyDef[*NumberDefinition, float64, *float64]
}

func Number() *NumberDefinition {
	var def NumberDefinition

	def.propertyDef = newPropertyDef[*NumberDefinition, float64](&def)

	return &def
}

func (nd *NumberDefinition) MarshalJSON() ([]byte, error) {
	if nd == nil {
		return []byte("null"), nil
	}

	var value struct {
		commonJSONSchema
		Enum []*float64 `json:"enum,omitempty"`
	}

	value.Type = "number"

	if nd.title.Set() {
		value.Title = ptr(nd.title.Value())
	}

	if nd.description.Set() {
		value.Description = ptr(nd.description.Value())
	}

	if nd.nullable.Set() {
		value.Nullable = ptr(nd.nullable.Value())
	}

	value.Enum = nd.enum

	bs, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode integer json: %w", err)
	}

	return bs, nil
}

type BooleanDefinition struct {
	propertyDef[*BooleanDefinition, bool, *bool]
}

func Boolean() *BooleanDefinition {
	var def BooleanDefinition

	def.propertyDef = newPropertyDef[*BooleanDefinition, bool](&def)

	return &def
}

func (bd *BooleanDefinition) MarshalJSON() ([]byte, error) {
	if bd == nil {
		return []byte("null"), nil
	}

	var value struct {
		commonJSONSchema
		Enum []*bool `json:"enum,omitempty"`
	}

	value.Type = "boolean"

	if bd.title.Set() {
		value.Title = ptr(bd.title.Value())
	}

	if bd.description.Set() {
		value.Description = ptr(bd.description.Value())
	}

	if bd.nullable.Set() {
		value.Nullable = ptr(bd.nullable.Value())
	}

	value.Enum = bd.enum

	bs, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode integer json: %w", err)
	}

	return bs, nil
}

type ArrayDefinition struct {
	propertyDef[*ArrayDefinition, struct{}, *struct{}]
	items    any
	minItems opt.T[uint]
	maxItems opt.T[uint]
}

func Array() *ArrayDefinition {
	var def ArrayDefinition

	def.propertyDef = newPropertyDef[*ArrayDefinition, struct{}](&def)

	return &def
}

func (ad *ArrayDefinition) Items(typDef any) *ArrayDefinition {
	ad.items = typDef
	return ad
}

func (ad *ArrayDefinition) MinItems(n uint) *ArrayDefinition {
	ad.minItems = opt.Wrap(n)
	return ad
}

func (ad *ArrayDefinition) MaxItems(n uint) *ArrayDefinition {
	ad.maxItems = opt.Wrap(n)
	return ad
}

func (ad *ArrayDefinition) MarshalJSON() ([]byte, error) {
	if ad == nil {
		return []byte("null"), nil
	}

	var value struct {
		commonJSONSchema
		Items    any   `json:"items,omitempty"`
		MinItems *uint `json:"minItems,omitempty"`
		MaxItems *uint `json:"maxItems,omitempty"`
	}

	value.Type = "array"

	if ad.title.Set() {
		value.Title = ptr(ad.title.Value())
	}

	if ad.description.Set() {
		value.Description = ptr(ad.description.Value())
	}

	if ad.nullable.Set() {
		value.Nullable = ptr(ad.nullable.Value())
	}

	if ad.enum != nil {
		return nil, errors.New("array enum is not implemented")
	}

	value.Items = ad.items

	if ad.minItems.Set() {
		value.MinItems = ptr(ad.minItems.Value())
	}

	if ad.maxItems.Set() {
		value.MaxItems = ptr(ad.maxItems.Value())
	}

	bs, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode integer json: %w", err)
	}

	return bs, nil
}

type ObjectDefinition struct {
	propertyDef[*ObjectDefinition, struct{}, *struct{}]

	propertiesNames []string
	properties      []any

	required             []string
	additionalProperties opt.T[bool]
}

func (od *ObjectDefinition) Required(properties ...string) *ObjectDefinition {
	od.required = append(od.required, properties...)
	return od
}

func (od *ObjectDefinition) Property(name string, prop any) *ObjectDefinition {
	od.propertiesNames = append(od.propertiesNames, name)
	od.properties = append(od.properties, prop)

	return od
}

func (od *ObjectDefinition) AdditionalProperties(p bool) *ObjectDefinition {
	od.additionalProperties = opt.Wrap(p)

	return od
}

func Object() *ObjectDefinition {
	var def ObjectDefinition

	def.propertyDef = newPropertyDef[*ObjectDefinition, struct{}](&def)

	return &def
}

func (od *ObjectDefinition) MarshalJSON() ([]byte, error) {
	if od == nil {
		return []byte("null"), nil
	}

	sort.Sort(&sortableProps{
		Names: od.propertiesNames,
		Props: od.properties,
	})

	for i, j := 0, -1; i < len(od.propertiesNames); i, j = i+1, j+1 {
		if j < 0 {
			continue
		}

		if l, r := od.propertiesNames[i], od.propertiesNames[j]; l == r {
			return nil, fmt.Errorf("duplicated definition of the property %q", l)
		}
	}

	for _, requiredName := range od.required {
		if i := sort.SearchStrings(od.propertiesNames, requiredName); i >= len(od.propertiesNames) || od.propertiesNames[i] != requiredName {
			return nil, fmt.Errorf("required property %q is not presented at the object definition", requiredName)
		}
	}

	var value struct {
		commonJSONSchema
		AdditionalProperties *bool                      `json:"additionalProperties,omitempty"`
		Required             []string                   `json:"required,omitempty"`
		Properties           map[string]json.RawMessage `json:"properties,omitempty"`
	}

	value.Type = "object"

	if od.title.Set() {
		value.Title = ptr(od.title.Value())
	}

	if od.description.Set() {
		value.Description = ptr(od.description.Value())
	}

	if od.additionalProperties.Set() {
		value.AdditionalProperties = ptr(od.additionalProperties.Value())
	}

	if od.nullable.Set() {
		value.Nullable = ptr(od.nullable.Value())
	}

	if len(od.propertiesNames) > 0 {
		value.Properties = make(map[string]json.RawMessage, len(od.propertiesNames))
	}

	if od.required != nil {
		value.Required = od.required
	}

	for i := range od.propertiesNames {
		name := od.propertiesNames[i]
		prop := od.properties[i]

		bs, err := json.Marshal(prop)
		if err != nil {
			return nil, fmt.Errorf("encode property %q as json: %w", name, err)
		}

		value.Properties[name] = bs
	}

	if od.enum != nil {
		return nil, errors.New("object enum is not implemented")
	}

	bs, err := json.Marshal(&value)
	if err != nil {
		return nil, fmt.Errorf("encode as json: %w", err)
	}

	return bs, nil
}

type property interface {
	*StringDefinition | *IntegerDefinition | *NumberDefinition | *BooleanDefinition | *ArrayDefinition | *ObjectDefinition
}

func ptr[T any](v T) *T {
	return &v
}
