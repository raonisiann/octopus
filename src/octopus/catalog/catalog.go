package catalog

// Catalog represents a group of resources
type Catalog struct {
    resources []Item
}

// Item defines the resource and it settings
type Item struct {
    resource    string
    class       string
    attributes  map[string]string
}

// New instantiate a new catalog
func New() Catalog {
    var catalog Catalog

    return catalog
}

// AddResource adds an item to catalog
func (c *Catalog) AddResource(i Item) {
    c.resources = append(c.resources, i)
}