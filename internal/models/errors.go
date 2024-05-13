package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")
var ErrTemplateCache = errors.New("template: newCache function error")
