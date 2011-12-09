package movo

import (
    html "fixedhtml"
    "strings"
    )

type NodeCursor struct {
  Node *html.Node
  parent *html.Node
  myIndex int
  Valid bool
}

// Node-like interface
func (c *NodeCursor) Parent() *NodeCursor {
  return NewCursor(c.parent)
}

func (c *NodeCursor) Child() []*NodeCursor {
  res := make([]*NodeCursor, len(c.Node.Child))

  for i, cd := range c.Node.Child {
    res[i] = NewCursor(cd)
  }
  return res
}

func (c *NodeCursor) Type() html.NodeType {
  return c.Node.Type
}

func (c *NodeCursor) Data() string {
  return c.Node.Data
}

func (c *NodeCursor) Attr() []html.Attribute {
  return c.Node.Attr
}

// Constructor
func NewCursor(n *html.Node) (res *NodeCursor) {
  res = &NodeCursor{}
  res.Node = n
  res.parent = n.Parent

  if res.parent != nil {
    for i, nn := range res.parent.Child {
      if nn == res.Node {
        res.myIndex = i
      }
    }
  } else {
    res.myIndex = -1
  }

  res.Valid = true

  return
}

// Navigation
func (c *NodeCursor) NextSibling() *NodeCursor {
  if c.parent != nil && len(c.parent.Child) > c.myIndex + 1 {
    return NewCursor(c.parent.Child[c.myIndex+1])
  }

  return InvalidCursor
}

// Search functions
func (c *NodeCursor) FindById(id string) (res *NodeCursor) {
  if !c.Valid {
    return c
  }

  for _,a := range c.Attr() {
    if a.Key == "id" && a.Val == id {
      return c
    }
  }

  for _,cc := range c.Child() {
    if res = cc.FindById(id); res.Valid {
      return res
    }
  }

  res = InvalidCursor

  return res
}

func (c *NodeCursor) FindText(text string) *NodeCursor {
  if !c.Valid {
    return c
  }

  if c.Data() == text {
    return c
  }

  for _,cc := range c.Child() {
    if res := cc.FindText(text); res.Valid {
      return res
    }
  }

  return InvalidCursor
}

func (c *NodeCursor) NextSiblingElement(elmType string) *NodeCursor {
  return c.findSiblingElement(elmType, 1)
}

func (c *NodeCursor) PrevSiblingElement(elmType string) *NodeCursor {
  return c.findSiblingElement(elmType, -1)
}

func (c *NodeCursor) findSiblingElement(elmType string, dir int8) (res *NodeCursor) {
  if !c.Valid {
    return c
  }

  for i := c.myIndex + 1; i < len(c.parent.Child); i+=int(dir) {
    cc := c.parent.Child[i]
    if cc.Type == html.ElementNode && cc.Data == elmType {
      res = &NodeCursor{Node: cc, parent: c.parent, myIndex: i, Valid: true}
      return res
    }
  }

  res = InvalidCursor

  return res
}

func (c *NodeCursor) FindChildElement(elmType string) (res *NodeCursor) {
  if !c.Valid {
    return c
  }

  for _,cc := range c.Child() {
    if cc.Type() == html.ElementNode && cc.Data() == elmType {
      return cc
    }
    res = cc.FindChildElement(elmType)
    if res.Valid {
      return res
    }
  }

  res = InvalidCursor

  return res
}

// Misc
func (c *NodeCursor) Prune() {
  for i := 0; c.Node.Child != nil && i < len(c.Node.Child); i++ {
    if cc := c.Node.Child[i]; cc.Type == html.TextNode {
      if cc.Data = strings.TrimSpace(cc.Data); len(cc.Data) == 0 {
        c.Node.Remove(cc)
        i--
      }
    } else if cc.Type == html.CommentNode || cc.Type == html.DoctypeNode {
      c.Node.Remove(cc)
      i--
    } else if cc.Type == html.ElementNode {
      if cc.Data == "br" || cc.Data == "input" || cc.Data == "hr" ||
          cc.Data == "head" {
        c.Node.Remove(cc)
        i--
      } else {
        NewCursor(cc).Prune()
      }
    }
  }
}

var (
    InvalidCursor *NodeCursor = &NodeCursor{}
    )

func init() {
  InvalidCursor.Valid = false
}
