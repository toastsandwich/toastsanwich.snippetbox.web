package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")
var ErrTemplateCache = errors.New("template: newCache function error")
var ErrDuplicateEmail = errors.New("models: duplicate email")
var ErrInvalidCredentials = errors.New("models: email or hash not found")
