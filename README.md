![flattenhtml](assets/flattenhtml.png)

flattenthtml is a Go package that helps you access to specific nodes in
a HTML document directly without a need for traversing all nodes.

![gerrors CI Flow](https://github.com/seinshah/flattenhtml/actions/workflows/ci.yaml/badge.svg)
[![Maintainability](https://api.codeclimate.com/v1/badges/13fba7eed22cc226da92/maintainability)](https://codeclimate.com/github/seinshah/flattenhtml/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/13fba7eed22cc226da92/test_coverage)](https://codeclimate.com/github/seinshah/flattenhtml/test_coverage)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/seinshah/flattenhtml?logo=github&sort=semver)

## Installation

```bash
go get github.com/seinshah/flattenhtml
```

## Overview

Use built-in or custom flatteners to access HTML document nodes directly
using your desired selectors. Whether you want to access all `div` nodes
(based on the tag name) or all elements with `class` attributes, or all
elements with `class` value as `container`, and so on.

`flattenhtml` currently supports the following flatteners out of the box:

- `TagFlattener`: flattens all nodes based on their tag name.

You can build a custom in-house flattener by implementing
`*flattenhtml.Flattener` interface. If your implementation is generic and
can be used by others, please consider contributing it to this package.

## Usage

```go
package main

import (
    "fmt"
    "log"
    "strings"

    "github.com/seinshah/flattenhtml"
)

func main() {
    // HTML document to be flattened.
    html := `
        <html>
            <head>
                <title>flattenhtml</title>
            </head>
            <body>
                <div class="container" id="target">
                    <div class="row">
                        <div class="col-md-6">
                            <h1>flattenhtml</h1>
                            <p>flattens HTML documents</p>
                        </div>
                        <div class="col-md-6">
                            <h1>flattenhtml</h1>
                            <p>flattens HTML documents</p>
                        </div>
                    </div>
                </div>
            </body>
        </html>
    `

    nm, err := flattenhtml.NewNodeManagerFromReader(strings.NewReader(html))
    if err != nil {
        log.Fatal(err)
    }

    mc, err := nm.Parse(flattenhtml.NewTagFlattener())
    if err != nil {
        log.Fatal(err)
    }

    tf, err := mc.SelectFlattener(&flattenhtml.TagFlattener{})
    if err != nil {
        log.Fatal(err)
    }

    divs := tf.SelectNodes("div")

    divs.
        Filter(flattenhtml.WithAttributeValueAs("class", "container")).
        Each(func(n *flattenhtml.Node) {
            val, _ := n.Attribute("id")

            fmt.Println(val)

            // Output:
            // target
        })
}
```
