package proto

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestParseFormattedProto2UnitTest(t *testing.T) {
	parseFormattedParsed(t, "unittest_proto2.proto")
}

func TestParseFormattedProto3UnitTest(t *testing.T) {
	parseFormattedParsed(t, "unittest_proto3.proto")
}

func TestParseFormattedProto3ArenaUnitTest(t *testing.T) {
	parseFormattedParsed(t, "unittest_proto3_arena.proto")
}

func parseFormattedParsed(t *testing.T, filename string) {
	// open it
	f, err := os.Open(filepath.Join("cmd", "protofmt", filename))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	// parse it
	p := NewParser(f)
	def, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	// count it
	c := new(counter)
	c.Count(def.Elements)
	beforeCount := c.count

	// format it
	out := new(bytes.Buffer)
	fmt := NewFormatter(out, "  ")
	fmt.Format(def)
	// parse the formatted content
	fp := NewParser(bytes.NewReader(out.Bytes()))
	_, err = fp.Parse()
	if err != nil {
		t.Fatal(err)
	}
	// count it again
	c.count = 0
	c.Count(def.Elements)
	afterCount := c.count
	if got, want := afterCount, beforeCount; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	t.Log("# proto elements", afterCount)
}

type counter struct {
	count int
}

func (c *counter) Count(elements []Visitee) {
	for _, each := range elements {
		each.Accept(c)
	}
}
func (c *counter) VisitMessage(m *Message) {
	c.count++
	c.Count(m.Elements)
}
func (c *counter) VisitService(v *Service) {
	c.count++
	c.Count(v.Elements)
}
func (c *counter) VisitSyntax(s *Syntax)           { c.count++ }
func (c *counter) VisitPackage(p *Package)         { c.count++ }
func (c *counter) VisitOption(o *Option)           { c.count++ }
func (c *counter) VisitImport(i *Import)           { c.count++ }
func (c *counter) VisitNormalField(i *NormalField) { c.count++ }
func (c *counter) VisitEnumField(i *EnumField)     { c.count++ }
func (c *counter) VisitEnum(e *Enum) {
	c.count++
	c.Count(e.Elements)
}
func (c *counter) VisitComment(e *Comment) { c.count++ }
func (c *counter) VisitOneof(o *Oneof) {
	c.count++
	c.Count(o.Elements)
}
func (c *counter) VisitOneofField(o *OneOfField) { c.count++ }
func (c *counter) VisitReserved(rs *Reserved)    { c.count++ }
func (c *counter) VisitRPC(rpc *RPC)             { c.count++ }
func (c *counter) VisitMapField(f *MapField)     { c.count++ }
func (c *counter) VisitGroup(g *Group) {
	c.count++
	c.Count(g.Elements)
}
func (c *counter) VisitExtensions(e *Extensions) { c.count++ }
