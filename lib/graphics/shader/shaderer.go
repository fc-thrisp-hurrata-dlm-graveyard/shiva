package shader

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/graphics/material"
	"github.com/Laughs-In-Flowers/shiva/lib/xrror"
)

type Shaderer interface {
	Templater
	//GetProg(string) *Prog
	//SetProg(string, string, string, string, string)
	GenerateProfile(material.Material) *Profile
	SetProgram(graphics.Provider, *Profile) error
	GenerateProgram(graphics.Provider, *Profile) (graphics.Program, error)
	Bind(graphics.Provider, *Bind) error
}

type shaderer struct {
	Templater
	prog []*Prog
	prgm []*Program
	curr *Program
}

func NewShaderer() *shaderer {
	return &shaderer{
		NewTemplater(NewLoaderSet(), EmptyFuncSet()),
		make([]*Prog, 0),
		make([]*Program, 0),
		nil,
	}
}

func DefaultShaderer() *shaderer {
	return &shaderer{
		defaultTemplater(),
		defaultProg,
		make([]*Program, 0),
		nil,
	}
}

func (s *shaderer) GetProg(tag string) *Prog {
	for _, p := range s.prog {
		if tag == p.Tag {
			return p
		}
	}
	return s.GetProg("basic")
}

func (s *shaderer) SetProg(tag, version, fragment, geometry, vertex string) {
	s.prog = append(s.prog, &Prog{tag, version, fragment, geometry, vertex})
}

func (s *shaderer) SetProgram(p graphics.Provider, pr *Profile) error {
	var h graphics.Program
	var err error
	for _, program := range s.prgm {
		if program.Equals(pr) {
			h = program.handle
			goto USEPROGRAM
		}
	}
	h, err = s.GenerateProgram(p, pr)
	goto USEPROGRAM
USEPROGRAM:
	p.UseProgram(h)
	return err
}

func (s *shaderer) GenerateProfile(m material.Material) *Profile {
	prog := s.GetProg(m.Shader())
	independent := m.Independent()
	var ambient, directional, point, spot int
	if !independent {
		//ambient
		//directional
		//point
		//spot
	}
	use := m.UseLights()
	matTexCt := m.TextureCount()
	return &Profile{
		prog,
		independent,
		use,
		ambient,
		directional,
		point,
		spot,
		matTexCt,
	}
}

func (s *shaderer) GenerateProgram(p graphics.Provider, pr *Profile) (graphics.Program, error) {
	prgm, err := NewProgram(p, pr, s)
	if err != nil {
		return 0, err
	}
	s.prgm = append(s.prgm, prgm)
	return prgm.handle, nil
}

func (s *shaderer) Bind(p graphics.Provider, b *Bind) error {
	return nil
}

type Templater interface {
	Render(io.Writer, string, interface{}) error
	Fetch(string) (*template.Template, error)
}

type templater struct {
	*LoaderSet
	*FuncSet
}

func NewTemplater(l *LoaderSet, f *FuncSet) *templater {
	return &templater{l, f}
}

func defaultTemplater() *templater {
	return NewTemplater(defaultLoaderSet(), defaultFuncSet())
}

func (t *templater) Render(w io.Writer, name string, data interface{}) error {
	tmpl, err := t.assemble(name)
	if err != nil {
		return err
	}

	if tmpl == nil {
		return NilTemplateError(name)
	}

	return tmpl.Execute(w, data)
}

func (t *templater) Fetch(name string) (*template.Template, error) {
	return t.assemble(name)
}

var (
	reExtendsTag  *regexp.Regexp = regexp.MustCompile("{{ extends [\"']?([^'\"}']*)[\"']? }}")
	reIncludeTag  *regexp.Regexp = regexp.MustCompile(`{{ include ["']?([^"]*)["']? }}`)
	reDefineTag   *regexp.Regexp = regexp.MustCompile("{{ ?define \"([^\"]*)\" ?\"?([a-zA-Z0-9]*)?\"? ?}}")
	reTemplateTag *regexp.Regexp = regexp.MustCompile("{{ ?template \"([^\"]*)\" ?([^ ]*)? ?}}")

	NilTemplateError   = xrror.Xrror("nil template named %s").Out
	NoTemplateError    = xrror.Xrror("no template named %s").Out
	EmptyTemplateError = xrror.Xrror("empty template named %s").Out
)

func (t *templater) assemble(name string) (*template.Template, error) {
	stack := make([]*Node, 0)

	err := t.add(&stack, name)

	if err != nil {
		return nil, err
	}

	blocks := map[string]string{}
	blockId := 0

	for _, node := range stack {
		var errInReplace error = nil
		node.Src = reIncludeTag.ReplaceAllStringFunc(node.Src, func(raw string) string {
			parsed := reIncludeTag.FindStringSubmatch(raw)
			templatePath := parsed[1]
			subTpl, err := t.getTemplate(templatePath)
			if err != nil {
				errInReplace = err
				return "[error]"
			}
			return subTpl
		})
		if errInReplace != nil {
			return nil, errInReplace
		}
	}

	for _, node := range stack {
		node.Src = reDefineTag.ReplaceAllStringFunc(node.Src, func(raw string) string {
			parsed := reDefineTag.FindStringSubmatch(raw)
			blockName := fmt.Sprintf("BLOCK_%d", blockId)
			blocks[parsed[1]] = blockName

			blockId += 1
			return "{{ define \"" + blockName + "\" }}"
		})
	}

	var rootTemplate *template.Template

	for i, node := range stack {
		node.Src = reTemplateTag.ReplaceAllStringFunc(node.Src, func(raw string) string {
			parsed := reTemplateTag.FindStringSubmatch(raw)
			origName := parsed[1]
			replacedName, ok := blocks[origName]
			dot := "."
			if len(parsed) == 3 && len(parsed[2]) > 0 {
				dot = parsed[2]
			}
			if ok {
				return fmt.Sprintf(`{{ template "%s" %s }}`, replacedName, dot)
			} else {
				return ""
			}
		})

		var thisTemplate *template.Template

		if i == 0 {
			thisTemplate = template.New(node.Name)
			rootTemplate = thisTemplate
		} else {
			thisTemplate = rootTemplate.New(node.Name)
		}

		thisTemplate.Funcs(t.GetFuncs())

		_, err := thisTemplate.Parse(node.Src)
		if err != nil {
			return nil, err
		}
	}

	return rootTemplate, nil
}

func (t *templater) getTemplate(name string) (string, error) {
	for _, l := range t.GetLoaders() {
		tmpl, err := l.Load(name)
		if err == nil {
			return tmpl, nil
		}
	}
	return "", NoTemplateError(name)
}

type Node struct {
	Name string
	Src  string
}

func (t *templater) add(stack *[]*Node, name string) error {
	tplSrc, err := t.getTemplate(name)
	if err != nil {
		return err
	}

	if len(tplSrc) < 1 {
		return EmptyTemplateError(name)
	}

	extendsMatches := reExtendsTag.FindStringSubmatch(tplSrc)
	if len(extendsMatches) == 2 {
		err := t.add(stack, extendsMatches[1])
		if err != nil {
			return err
		}
		tplSrc = reExtendsTag.ReplaceAllString(tplSrc, "")
	}

	node := &Node{
		Name: name,
		Src:  tplSrc,
	}

	*stack = append((*stack), node)

	return nil
}

type FuncSet struct {
	f map[string]interface{}
}

func NewFuncSet(fns map[string]interface{}) *FuncSet {
	return &FuncSet{fns}
}

func EmptyFuncSet() *FuncSet {
	return &FuncSet{make(map[string]interface{})}
}

var defaultFuncs = map[string]interface{}{
	"loop": func(n int) []int {
		s := make([]int, n)
		for i := range s {
			s[i] = i
		}
		return s
	},
}

func defaultFuncSet() *FuncSet {
	return NewFuncSet(defaultFuncs)
}

func (f *FuncSet) AddFuncs(fns map[string]interface{}) {
	for k, fn := range fns {
		f.f[k] = fn
	}
}

func (f *FuncSet) GetFuncs() map[string]interface{} {
	return f.f
}

type LoaderSet struct {
	l []Loader
}

func NewLoaderSet() *LoaderSet {
	return &LoaderSet{make([]Loader, 0)}
}

func defaultLoaderSet() *LoaderSet {
	ls := NewLoaderSet()
	ls.AddLoaders(
		FileLoader(),
		MapLoader(defaultTemplates),
	)
	return ls
}

func (l *LoaderSet) AddLoaders(ls ...Loader) {
	l.l = append(l.l, ls...)
}

func (l *LoaderSet) GetLoaders(ls ...Loader) []Loader {
	return l.l
}

type Loader interface {
	Load(string) (string, error)
	ListTemplates() ([]string, error)
}

type BaseLoader struct {
	Errors         []error
	FileExtensions []string
}

var NoMethodError = xrror.Xrror("%s method not implemented").Out

func (b *BaseLoader) Load(name string) (string, error) {
	return "", NoMethodError("load")
}

func (b *BaseLoader) ListTemplates() ([]string, error) {
	return []string{}, NoMethodError("list templates")
}

func (b *BaseLoader) ValidExtension(ext string) bool {
	for _, extension := range b.FileExtensions {
		if extension == ext {
			return true
		}
	}
	return false
}

var PathError = xrror.Xrror("path: %s returned error").Out

type fileLoader struct {
	BaseLoader
}

func FileLoader() *fileLoader {
	return &fileLoader{}
}

func (f *fileLoader) Load(path string) (string, error) {
	if f.ValidExtension(filepath.Ext(path)) {
		if _, err := os.Stat(path); err == nil {
			file, err := os.Open(path)
			r, err := ioutil.ReadAll(file)
			return string(r), err
		}
	}
	return "", NoTemplateError(path)
}

type mapLoader struct {
	BaseLoader
	TemplateMap map[string]string
}

func MapLoader(tm ...map[string]string) *mapLoader {
	m := &mapLoader{TemplateMap: make(map[string]string)}
	for _, t := range tm {
		for k, v := range t {
			m.TemplateMap[k] = v
		}
	}
	return m
}

func (l *mapLoader) Load(name string) (string, error) {
	if r, ok := l.TemplateMap[name]; ok {
		return string(r), nil
	}
	return "", NoTemplateError(name)
}

func (l *mapLoader) ListTemplates() ([]string, error) {
	var listing []string
	for k, _ := range l.TemplateMap {
		listing = append(listing, k)
	}
	return listing, nil
}
