package globaloption

import (
	"errors"
	"strings"
)

const OptionNameAuthors string = "authors"

type Author string

func parseAuthor(v string) (Author, error) {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return "", errors.New("empty author is not allowed")
	}
	return Author(trimmed), nil
}

func (a Author) String() string {
	return string(a)
}

type Authors []Author

func ParseAuthors(v string) (Authors, error) {
	v = strings.Trim(v, ", ")
	if v == "" {
		return Authors{}, nil
	}

	values := strings.Split(v, ",")
	authors := make([]Author, 0, len(values))
	for _, value := range values {
		author, err := parseAuthor(value)
		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (a Authors) Strings() []string {
	result := make([]string, 0, len(a))
	for _, author := range a {
		result = append(result, author.String())
	}
	return result
}
