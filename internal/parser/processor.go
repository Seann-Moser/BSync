package parser

import (
	"encoding/json"
	v2 "github.com/Seann-Moser/WebParser/v2"
	"github.com/Seann-Moser/WebParser/website"
	"go.uber.org/zap"
	"net/url"
)

type ParserProcessor struct {
	Req    *v2.HTMLSourceRequest
	Logger *zap.Logger
}

func NewParserProcessor(Logger *zap.Logger) *ParserProcessor {
	return &ParserProcessor{
		Req:    v2.NewHTMLSourceRequestWithSleep(5),
		Logger: Logger,
	}
}

func (pp *ParserProcessor) GetData(u string, searchData []*website.Search) ([]byte, error) {
	if len(searchData) == 0 {
		return nil, nil
	}
	newURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	p, err := website.NewWebParser(newURL.Host, u, searchData)
	if err != nil {
		return nil, err
	}
	_, rmap, err := p.Parse(pp.Req, u)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(rmap, "", "\t")
}
