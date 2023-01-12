package etymology

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

func parseHTML(htmlStr string) (map[string]string, error) {
	rawHTML := string(htmlStr)
	tokenizer := html.NewTokenizer(strings.NewReader(rawHTML))
	output := make(map[string]string)
	for {
		//get the next token type
		tokenType := tokenizer.Next()

		if eof, err := endOfParsing(tokenizer, tokenType); err != nil {
			return nil, err
		} else if eof {
			// EOF, no more parsing to be done
			break
		}

		//process the token according to the token type...
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			className, exists := getAttr(token, "class")
			if !exists {
				continue
			}
			isWord := className[:6] == "word--"
			if token.Data == "div" && isWord {
				name, history, eof, err := extractHistory(tokenizer)
				if err != nil {
					return nil, err
				}
				if eof {
					break
				}
				output[name] = history
				if len(output) >= maxDefinitions {
					break
				}
			}
		}
	}
	return output, nil
}

func extractHistory(tokenizer *html.Tokenizer) (name string, definition string, eof bool, err error) {
	// get word with type
	for name == "" || definition == "" {
		tokenType := tokenizer.Next()
		if eof, err := endOfParsing(tokenizer, tokenType); err != nil {
			return "", "", eof, err
		} else if eof {
			// EOF, no more parsing to be done
			break
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			className, exists := getAttr(token, "class")

			if token.Data == "a" && exists && className[:12] == "word__name--" {
				name, eof, err = extractTitle(tokenizer)
				if eof || err != nil {
					return "", "", eof, err
				}
			}
			if token.Data == "section" && exists && className[:18] == "word__defination--" {
				definition, eof, err = extractDefinition(tokenizer)
				if eof || err != nil {
					return "", "", eof, err
				}
			}
		}
	}
	return name, definition, false, nil
}

func extractTitle(tokenizer *html.Tokenizer) (title string, eof bool, err error) {
	tokenType := tokenizer.Next()
	token := tokenizer.Token()
	for tokenType != html.EndTagToken && token.Data != "a" {
		if eof, err := endOfParsing(tokenizer, tokenType); eof || err != nil {
			return "", false, err
		}

		if tokenType == html.TextToken {
			title += token.Data
		}
		tokenType = tokenizer.Next()
		token = tokenizer.Token()
	}
	return title, eof, err
}

func extractDefinition(tokenizer *html.Tokenizer) (definition string, eof bool, err error) {
	// TODO: Re-add formatting (ie bold, italics) for the tags
	tokenType := tokenizer.Next()
	token := tokenizer.Token()
	for tokenType != html.EndTagToken || token.Data != "section" {
		if eof, err := endOfParsing(tokenizer, tokenType); eof || err != nil {
			break
		}

		if tokenType == html.TextToken {
			definition += token.Data
		} else if tokenType == html.StartTagToken && token.Data == "blockquote" {
			definition += "Quote:\n"
		} else if tokenType == html.EndTagToken && token.Data == "blockquote" {
			definition += "\n\n"
		} else if tokenType == html.EndTagToken && token.Data == "p" {
			definition += "\n\n"
		}
		tokenType = tokenizer.Next()
		token = tokenizer.Token()
	}

	definition = strings.ReplaceAll(definition, "\u00a0", " ")
	definition = strings.ReplaceAll(definition, "*", "\\*")
	definition = strings.ReplaceAll(definition, "_", "\\_")
	return definition, eof, err
}
func getAttr(token html.Token, key string) (string, bool) {
	for _, attr := range token.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func endOfParsing(tokenizer *html.Tokenizer, tokenType html.TokenType) (eof bool, err error) {
	if tokenType == html.ErrorToken {
		err := tokenizer.Err()
		if err == io.EOF {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
