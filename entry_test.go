package miniplan

import "fmt"

func ExampleEntry_Tags() {
	e := Entry{
		Title:       "Color the wall #kitchen, #house",
		Description: "Something warm #color",
	}
	fmt.Println(e.Tags())
	// output:
	// [#kitchen #house #color]
}
