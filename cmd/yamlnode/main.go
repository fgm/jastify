package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func alias() *yaml.Node {
	const anchor = "ref"
	sn := &yaml.Node{
		Kind: yaml.ScalarNode, Value: "key", Anchor: anchor,
		HeadComment: "HC Key",
		FootComment: "FC Key",
		LineComment: "LC Key",
	}
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			&yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					// Notice how LineComment is emitted on the NEXT value for map keys and FootComment is emitted after current value
					sn,
					// Note how value HeadComment is emitted AFTER the key FootComment
					{Kind: yaml.ScalarNode, Value: "a value", HeadComment: "HC Value", FootComment: "FC2 Value", LineComment: "LC Value"},

					{Kind: yaml.ScalarNode, Value: "a key for int"}, // See: no LC, but we get the one for the previous key
					{Kind: yaml.ScalarNode, Tag: "!!int", Value: "12"},

					{Kind: yaml.ScalarNode, Value: "a key for null"}, // See: no LC, but we get the one for the previous key
					{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"},

					{Kind: yaml.ScalarNode, Value: "key to alias"},
					{Kind: yaml.AliasNode, Value: anchor, Alias: sn},
				},
			},
		},
	}
}

func listDefault() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.SequenceNode,
				Tag:  "!!seq",
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Value: "a",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "b",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "c",
					},
				},
			},
		},
	}
}

func mapDefault() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Tag:  "!!map",
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "k",
					},
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "v",
					},
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "l",
					},
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "w",
					},
				},
			},
		},
	}
}

func stringFlow() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind:        yaml.ScalarNode,
				Style:       yaml.FlowStyle,
				Tag:         "!!str",
				Value:       "a value",
				HeadComment: "Head",
				LineComment: "Line",
				FootComment: "Foot",
			},
		},
	}
}

func add(examples map[string]*yaml.Node, list *[]string, label string, gen func() *yaml.Node) {
	examples[label] = gen()
	*list = append(*list, label)
}

func main() {
	enc := yaml.NewEncoder(os.Stdout)

	var docAndMapping yaml.Node
	if err := yaml.Unmarshal([]byte(`
&ref k: a value
l: *ref
`), &docAndMapping); err != nil {
		log.Fatal(err)
	}
	//spew.Dump(docAndMapping.Content[0])

	examples := map[string]*yaml.Node{}
	keys := []string{}

	// This is wrong, but works
	//add(examples, &keys, "raw scalar", rawScalar)
	//add(examples, &keys, "string, tagged", stringT)
	//add(examples, &keys, "string, double quoted", stringDQ)
	//add(examples, &keys, "string, single quoted", stringSQ)
	//add(examples, &keys, "string, literal", stringLiteral)
	//add(examples, &keys, "string, folded", stringFolded)
	add(examples, &keys, "string, flow", stringFlow)
	add(examples, &keys, "map, default", mapDefault)
	add(examples, &keys, "list, default", listDefault)
	add(examples, &keys, "alias, default", alias)

	for _, key := range keys {
		fmt.Fprintln(os.Stdout, key)
		if err := enc.Encode(examples[key]); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(os.Stdout)
	}
}
