package main

import (
	"bytes"
	"strings"
	"path/filepath"
	"github.com/russross/blackfriday"
)
 
type ConfluenceRenderer struct {
	blackfriday.Renderer
	basePath string
	Images * map[string]MacroImage
}

func (renderer ConfluenceRenderer) BlockCode(
	out *bytes.Buffer,
	text []byte,
	lang string,
) {
	out.WriteString(MacroCode{lang, text}.Render())
}


func doubleSpace(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}



// Using if statements is a bit faster than a switch statement. As the compiler
// improves, this should be unnecessary this is only worthwhile because
// attrEscape is the single largest CPU user in normal use.
// Also tried using map, but that gave a ~3x slowdown.
func escapeSingleChar(char byte) (string, bool) {
	if char == '"' {
		return "&quot;", true
	}
	if char == '&' {
		return "&amp;", true
	}
	if char == '<' {
		return "&lt;", true
	}
	if char == '>' {
		return "&gt;", true
	}
	return "", false
}

func attrEscape(out *bytes.Buffer, src []byte) {
	org := 0
	for i, ch := range src {
		if entity, ok := escapeSingleChar(ch); ok {
			if i > org {
				// copy all the normal characters since the last escape
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString(entity)
		}
	}
	if org < len(src) {
		out.Write(src[org:])
	}
}

func (renderer ConfluenceRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	logger.Tracef("AutoLink: %s", link)

	options, ok := renderer.Renderer.(*blackfriday.Html)
	if ok {
		if  !strings.Contains(string(link),  "://") && strings.HasSuffix(string(link), ".md") {
			link = link[:len(link)-3]
		} 
		options.AutoLink(out,link,kind)	
	}
	
}

func (renderer ConfluenceRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {

	logger.Tracef("link: %s", link)

	options, ok := renderer.Renderer.(*blackfriday.Html)
	if ok {
		strLink := string(link)
		if  !strings.Contains(strLink,  "://") && strings.HasSuffix(strLink, ".md") {
			strLink =  ( filepath.Base(strLink) )
			link =  []byte (strLink)
			link = link[:len(link)-3]
			logger.Tracef("link: %s", link)
		} 
		options.Link(out,link,title,content)	
	}
}


func (renderer ConfluenceRenderer) Image (
	out *bytes.Buffer, 
	link []byte, 
	title []byte, 
	alt []byte,
) {
	if !strings.Contains(string(link), "://") {
	 	existing_macro, ok := (*renderer.Images)[string(link)]
		if ok {
			existing_macro.Render()
			return
		}
		logger.Tracef("Local Image: %s (%s) -> %s\n", string(title), string(alt), string(link))
		macro, err := newMacroImage(filepath.Join(renderer.basePath,  string(link)), string(title), string(alt) )
		if (err == nil) {
			(*renderer.Images)[string(link)] = *macro
			out.WriteString(macro.Render())
			return;
		} else {
			logger.Warningf("Local Image: %s (%s) -> %s was not uploaded (%v) \n", string(title), string(alt), string(link), err)
		}
	} 
	
	out.WriteString("<img alt=\"")
	if len(alt) > 0 {
		attrEscape(out, alt)
	}
	out.WriteString("\"")
	if len(title) > 0 {
		out.WriteString(" title=\"")
		attrEscape(out, title)
		out.WriteString("\"")
	}
	out.WriteString(" src=\"")
	// options.maybeWriteAbsolutePrefix(out, link)
	attrEscape(out, link)

	out.WriteString("\" />")
}