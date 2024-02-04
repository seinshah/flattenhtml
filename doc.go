// Package flattenhtml provides a way to flatten the HTML tree structure and
// then use the flattened data to do different kinds of lookups.
//
// Go provides [html] package that bear the heavy load of parsing HTML.
// However, this package results in a tree structure. Although it is generic
// and can be utilized for any traversal purposes, it is not very convenient
// for some use cases, such as, searching for a specific element.
//
// Here is where flattenhtml comes in. It provides different mechanism to flatten
// the HTML tree structure based on the use case. For example, if you want to
// work with the nodes based on their tag name, you can use TagFlattener flattener
// to first flatten all the nodes based on their tag name and then do continues tag
// lookup without the need for constantly traversing the tree.
//
// TagFlattener is currently the only built-in flattener of this package. However,
// all flatteners implement flattenhtml.Flattener interface and you can easily
// implement your own flattener.
//
// When you use the following statement to initialize the NodeManager, parsed HTML
// tree will be traversed once and for any further lookups, the flattener data is
// accessible without the need for traversing the tree again. Also, there is the
// possibility of using multiple flatteners at the same time. For example, you can
// use TagFlattener to flatten the nodes based on their tag name and then use
// AttributeFlattener to flatten the nodes based on their attributes. The same as
// before, the HTML tree will be traversed only once to utilize all flatteners.
//
//	html := "<html><head></head><body><div><p></p></div></body></html>"
//	flatteners := []flattenhtml.Flattener{flattenhtml.TagFlattener, ...}
//	nm := flattenhtml.NewNodeManagerFromReader(strings.NewReader(html))
//	mc := nm.Parse(flatteners...)
//
// Once the flattening process is done, you will have a *flattenhtml.MultiCursor
// Which holds a pointer to all the flatteners. Now, before proceeding, you need
// to select a single flattener of your choice, to continue the lookup process.
//
//	tagFlattenerCursor := mc.First()
//
// Now, you can get nodes of the same tag name using the following statement:
//
//	nodes := tagFlattenerCursor.SelectNodes("div")
//
// This will return a *flattenhtml.NodeIterator that can be used to iterate over
// the nodes that are selected by the given key. In this case, all the nodes that
// have "div" tag name.
//
// Note that the underlying engine for parsing the HTML is [golang.org/x/net/html]
// package and all the fact about standardizing the HTML tree applies to this package.
//
// [html]: https://pkg.go.dev/golang.org/x/net/html
// [golang.org/x/net/html]: https://pkg.go.dev/golang.org/x/net/html
package flattenhtml
