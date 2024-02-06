package reader

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type recipes struct {
	Cake []struct {
		Name        string `xml:"name" json:"name"`
		Time        string `xml:"stovetime" json:"time"`
		Ingredients []struct {
			IngredientName  string `xml:"itemname" json:"ingredient_name"`
			IngredientCount string `xml:"itemcount" json:"ingredient_count"`
			IngredientUnit  string `xml:"itemunit,omitempty" json:"ingredient_unit,omitempty"`
		} `xml:"ingredients>item" json:"ingredients"`
	} `xml:"cake" json:"cake"`
}

type DBReader interface {
	Read(s *string)
	Convert() string
	Compare(DBReader)
}

type Xml struct {
	data recipes
}

func newXml() *Xml {
	return &Xml{
		data: recipes{},
	}
}

func (x *Xml) Read(s *string) {
	txt, err := os.ReadFile(*s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
		os.Exit(1)
	}
	if err = xml.Unmarshal(txt, &x.data); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка декодирования:", err)
		os.Exit(1)
	}
}

func (x *Xml) Convert() string {
	xmlBytes, err := json.MarshalIndent(x.data, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка сериализации:", err)
		return ""

	}
	return string(xmlBytes)
}

type Json struct {
	data recipes
}

func newJson() *Json {
	return &Json{
		data: recipes{},
	}
}

func (j *Json) Read(s *string) {
	txt, err := os.ReadFile(*s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
		os.Exit(1)
	}
	if err = json.Unmarshal(txt, &j.data); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка декодирования:", err)
		os.Exit(1)
	}
}

func (j *Json) Convert() string {
	jsonBytes, err := xml.MarshalIndent(j.data, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка сериализации:", err)
		return ""
	}
	return string(jsonBytes)
}

func NewDBReader(s *string) (DBReader, error) {
	fileExt := filepath.Ext(*s)
	if fileExt == ".xml" {
		return newXml(), nil
	} else if fileExt == ".json" {
		return newJson(), nil
	} else {
		return nil, fmt.Errorf("Неподдерживаемый тип файла")
	}
}

func (x *Xml) Compare(j DBReader) {
	switch o := j.(type) {
	case *Xml:
		checkChangesDifferences(x.data, o.data)
	case *Json:
		checkChangesDifferences(x.data, o.data)
	}
}

func (j *Json) Compare(x DBReader) {
	switch o := x.(type) {
	case *Xml:
		checkChangesDifferences(j.data, o.data)
	case *Json:
		checkChangesDifferences(j.data, o.data)
	}
}

func checkChangesDifferences(old, new recipes) {
	indexesRecipes := checkCake(old, new)
	checkTime(old, new, indexesRecipes)
	checkIngredient(old, new, indexesRecipes)
	checkChangeIngredients(old, new, indexesRecipes)
}

func checkCake(old, new recipes) map[int]int {
	indexesRecipes := make(map[int]int)
	oldCake := make(map[string]int)
	for i, cake := range old.Cake {
		oldCake[cake.Name] = i
	}
	for indexNew, cake := range new.Cake {
		name := cake.Name
		indexOld, ok := oldCake[name]
		delete(oldCake, name)
		if !ok {
			fmt.Printf("ADDED cake \"%s\"\n", name)
		} else {
			indexesRecipes[indexOld] = indexNew
		}
	}
	for name := range oldCake {
		fmt.Printf("REMOVED cake \"%s\"\n", name)
	}
	return indexesRecipes
}

func checkTime(old, new recipes, indexesRecipes map[int]int) {
	for indexOld, indexNew := range indexesRecipes {
		if old.Cake[indexOld].Time != new.Cake[indexNew].Time {
			fmt.Printf(
				"CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n",
				old.Cake[indexOld].Name, new.Cake[indexNew].Time, old.Cake[indexOld].Time,
			)
		}
	}
}

func checkIngredient(old, new recipes, indexesRecipes map[int]int) {
	for indexOld, indexNew := range indexesRecipes {
		oldIngredient := make(map[string]int)
		for i, ingredient := range old.Cake[indexOld].Ingredients {
			oldIngredient[ingredient.IngredientName] = i
		}
		newIngredient := make(map[string]int)
		for i, ingredient := range new.Cake[indexNew].Ingredients {
			newIngredient[ingredient.IngredientName] = i
		}

		for ingredientName := range newIngredient {
			_, ok := oldIngredient[ingredientName]
			if !ok {
				fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", ingredientName, new.Cake[indexNew].Name)
			}
		}
		for ingredientName := range oldIngredient {
			_, ok := newIngredient[ingredientName]
			if !ok {
				fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", ingredientName, old.Cake[indexOld].Name)
			}
		}

	}
}

func checkChangeIngredients(old, new recipes, indexesRecipes map[int]int) {
	for indexOld, indexNew := range indexesRecipes {
		for _, oldIngredient := range old.Cake[indexOld].Ingredients {
			for _, newIngredient := range new.Cake[indexNew].Ingredients {
				if oldIngredient.IngredientName == newIngredient.IngredientName {
					if oldIngredient.IngredientCount != newIngredient.IngredientCount {
						fmt.Printf(
							"CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\" \n",
							newIngredient.IngredientName, new.Cake[indexNew].Name, newIngredient.IngredientCount,
							oldIngredient.IngredientCount,
						)
					}
					if oldIngredient.IngredientUnit != newIngredient.IngredientUnit {
						if oldIngredient.IngredientUnit == "" {
							fmt.Printf(
								"ADDED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n",
								newIngredient.IngredientUnit, oldIngredient.IngredientName, old.Cake[indexOld].Name,
							)
						} else if newIngredient.IngredientUnit == "" {
							fmt.Printf(
								"REMOVED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n",
								oldIngredient.IngredientUnit, oldIngredient.IngredientName, old.Cake[indexOld].Name,
							)
						} else {
							fmt.Printf(
								"CHANGED unit for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\" \n",
								oldIngredient.IngredientName, old.Cake[indexOld].Name, newIngredient.IngredientUnit,
								oldIngredient.IngredientUnit,
							)
						}
					}
				}
			}
		}
	}
}
