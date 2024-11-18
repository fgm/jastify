package helpers

import "gopkg.in/yaml.v3"

func RawScalar() *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: "a value",
	}
}

func StringDQ() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind:        yaml.ScalarNode,
				Style:       yaml.DoubleQuotedStyle,
				Tag:         "!!str",
				Value:       "a value",
				HeadComment: "Head",
				LineComment: "Line",
				FootComment: "Foot",
			},
		},
	}
}

func StringSQ() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind:        yaml.ScalarNode,
				Style:       yaml.SingleQuotedStyle,
				Tag:         "!!str",
				Value:       "a value",
				HeadComment: "Head",
				LineComment: "Line",
				FootComment: "Foot",
			},
		},
	}
}

func StringT() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind:        yaml.ScalarNode,
				Style:       yaml.TaggedStyle,
				Tag:         "!!str",
				Value:       "a value",
				HeadComment: "Head",
				LineComment: "Line",
				FootComment: "Foot",
			},
		},
	}
}

func StringLiteral() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind:        yaml.ScalarNode,
				Style:       yaml.LiteralStyle,
				Tag:         "!!str",
				Value:       "a value",
				HeadComment: "Head",
				LineComment: "Line",
				FootComment: "Foot",
			},
		},
	}
}

func StringFolded() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind:        yaml.ScalarNode,
				Style:       yaml.FoldedStyle,
				Tag:         "!!str",
				Value:       "a value",
				HeadComment: "Head",
				LineComment: "Line",
				FootComment: "Foot",
			},
		},
	}
}
