package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

const (
	manaShift   = 0
	healthShift = 10
	houseShift  = 20
	gunShift    = 21
	familyShift = 22
	typeShift   = 23

	manaMask   = 0x3FF
	healthMask = 0x3FF
	flagMask   = 0x1
	typeMask   = 0x3
)

const (
	respectShift    = 0
	strengthShift   = 4
	experienceShift = 8
	levelShift      = 12

	statMask = 0xF
)

func WithName(name string) Option {
	return func(p *GamePerson) {
		n := len(name)
		if n > len(p.name) {
			n = len(p.name)
		}
		copy(p.name[:], name[:n])
	}
}

func WithCoordinates(x, y, z int) Option {
	return func(p *GamePerson) {
		p.x = int32(x)
		p.y = int32(y)
		p.z = int32(z)
	}
}

func WithGold(gold int) Option {
	return func(p *GamePerson) {
		p.gold = uint32(gold)
	}
}

func WithMana(mana int) Option {
	return func(p *GamePerson) {
		p.attribute &^= manaMask << manaShift
		p.attribute |= (uint32(mana) & manaMask) << manaShift
	}
}

func WithHealth(health int) Option {
	return func(p *GamePerson) {
		p.attribute &^= healthMask << healthShift
		p.attribute |= (uint32(health) & healthMask) << healthShift
	}
}

func WithHouse() Option {
	return func(p *GamePerson) {
		p.attribute &^= flagMask << houseShift
		p.attribute |= flagMask << houseShift
	}
}

func WithGun() Option {
	return func(p *GamePerson) {
		p.attribute &^= flagMask << gunShift
		p.attribute |= flagMask << gunShift
	}
}

func WithFamily() Option {
	return func(p *GamePerson) {
		p.attribute &^= flagMask << familyShift
		p.attribute |= flagMask << familyShift
	}
}

func WithType(personType int) Option {
	return func(p *GamePerson) {
		p.attribute &^= typeMask << typeShift
		p.attribute |= (uint32(personType) & typeMask) << typeShift
	}
}

func WithRespect(respect int) Option {
	return func(p *GamePerson) {
		p.stats &^= statMask << respectShift
		p.stats |= (uint16(respect) & statMask) << respectShift
	}
}

func WithStrength(strength int) Option {
	return func(p *GamePerson) {
		p.stats &^= statMask << strengthShift
		p.stats |= (uint16(strength) & statMask) << strengthShift
	}
}

func WithExperience(experience int) Option {
	return func(p *GamePerson) {
		p.stats &^= statMask << experienceShift
		p.stats |= (uint16(experience) & statMask) << experienceShift
	}
}

func WithLevel(level int) Option {
	return func(p *GamePerson) {
		p.stats &^= statMask << levelShift
		p.stats |= (uint16(level) & statMask) << levelShift
	}
}

func NewGamePerson(options ...Option) GamePerson {
	var p GamePerson
	for _, opt := range options {
		opt(&p)
	}

	return p
}

func (p *GamePerson) Name() string {
	end := len(p.name)
	for i, b := range p.name {
		if b == 0 {
			end = i
			break
		}
	}

	return string(p.name[:end])
}

func (p *GamePerson) X() int          { return int(p.x) }
func (p *GamePerson) Y() int          { return int(p.y) }
func (p *GamePerson) Z() int          { return int(p.z) }
func (p *GamePerson) Gold() int       { return int(p.gold) }
func (p *GamePerson) Mana() int       { return int((p.attribute >> manaShift) & manaMask) }
func (p *GamePerson) Health() int     { return int((p.attribute >> healthShift) & healthMask) }
func (p *GamePerson) HasHouse() bool  { return ((p.attribute >> houseShift) & flagMask) == 1 }
func (p *GamePerson) HasGun() bool    { return ((p.attribute >> gunShift) & flagMask) == 1 }
func (p *GamePerson) HasFamily() bool { return ((p.attribute >> familyShift) & flagMask) == 1 }
func (p *GamePerson) Type() int       { return int((p.attribute >> typeShift) & typeMask) }
func (p *GamePerson) Respect() int    { return int((p.stats >> respectShift) & statMask) }
func (p *GamePerson) Strength() int   { return int((p.stats >> strengthShift) & statMask) }
func (p *GamePerson) Experience() int { return int((p.stats >> experienceShift) & statMask) }
func (p *GamePerson) Level() int      { return int((p.stats >> levelShift) & statMask) }

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z   int32    // координаты по 4 байта = 12 байт
	gold      uint32   // золото 4 байта
	attribute uint32   // 4 байта(0-9 бит мана, 10-19 здоровье, 20 наличие дома, 21 наличие оружия, 22 наличие семьи, 23-24 тип игрока)
	stats     uint16   // 2 байта(0-3 бит уважение, 4-7 сила, 8-11 опыт, 12 - 15 уровень)
	name      [42]byte // 42 байта
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
