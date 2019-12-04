package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var rGroup = regexp.MustCompile(`(\d+) units each with (\d+) hit points (?:\((weak|immune) to ([a-z, ]+)(?:; (weak|immune) to ([a-z, ]+))?\) )?with an attack that does (\d+) (\w+) damage at initiative (\d+)`)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath is passed")
	}
	fileName := os.Args[1]
	if immune, infections, err := getGroups(fileName); err != nil {
		log.Fatal(err)
	} else {
		res, _ := fight(copyGroups(immune), copyGroups(infections), 0)
		fmt.Println("Part1 result is", res)
		res2 := fightUntilVictory(immune, infections)
		fmt.Println("Part2 result is", res2)
	}
}

type GroupType int

const (
	Immune GroupType = iota
	Infection
)

func (g GroupType) String() string {
	if g == Immune {
		return "IS"
	} else {
		return "IN"
	}
}

type Group struct {
	ID            int
	Units         int
	HP            int
	Weaknesses    map[string]bool
	Immunities    map[string]bool
	AttackDammage int
	AttackType    string
	Initiative    int
	Type          GroupType
}

func (g Group) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("(%s%d) %d units | %d HPs | %d %s | I%d", g.Type, g.ID, g.Units, g.HP, g.AttackDammage, g.AttackType, g.Initiative))
	if len(g.Weaknesses) > 0 {
		s.WriteString(" | Weak to ")
		for v := range g.Weaknesses {
			s.WriteString(v + " ")
		}
	}
	if len(g.Immunities) > 0 {
		s.WriteString(" | Immune to ")
		for v := range g.Immunities {
			s.WriteString(v + " ")
		}
	}
	return s.String()
}

func (g Group) EP() int {
	return g.Units * g.AttackDammage
}

func (g Group) Dammage(gg Group) int {
	for c := range gg.Immunities {
		if c == g.AttackType {
			return 0
		}
	}
	for c := range gg.Weaknesses {
		if c == g.AttackType {
			return g.EP() * 2
		}
	}
	return g.EP()
}

func getGroups(fileName string) ([]Group, []Group, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var immune, infections []Group
	var curentType GroupType
	i := 1
	for scanner.Scan() {
		line := scanner.Text()
		if line == "Infection:" {
			curentType = Infection
			i = 1
		} else if !(line == "" || line == "Immune System:") {
			extracts := rGroup.FindAllStringSubmatch(line, -1)
			if len(extracts) != 1 || len(extracts[0]) != 10 {
				return nil, nil, errors.New("Can't parse instruction line " + line)
			}
			extract := extracts[0]
			weaknesses := make(map[string]bool)
			immunities := make(map[string]bool)
			group := Group{
				ID:            i,
				Units:         atoi(extract[1]),
				HP:            atoi(extract[2]),
				AttackDammage: atoi(extract[7]),
				AttackType:    extract[8],
				Initiative:    atoi(extract[9]),
				Weaknesses:    weaknesses,
				Immunities:    immunities,
				Type:          curentType,
			}

			setWeaknessesImmunities(extract[3], extract[4], weaknesses, immunities)
			setWeaknessesImmunities(extract[5], extract[6], weaknesses, immunities)
			if curentType == Immune {
				immune = append(immune, group)
			} else {
				infections = append(infections, group)
			}
			i++
		}
	}
	return immune, infections, nil
}

func setWeaknessesImmunities(level, list string, weaknesses, immunities map[string]bool) {
	conditionList := strings.Split(list, ", ")
	var levelToSet map[string]bool
	switch level {
	case "immune":
		levelToSet = immunities
	case "weak":
		levelToSet = weaknesses
	}
	if levelToSet != nil {
		for _, condition := range conditionList {
			levelToSet[condition] = true
		}
	}
}

func printGroups(groups []*Group) {
	for _, g := range groups {
		fmt.Println(g)
	}
}

func copyGroups(groups []Group) (newGroups []*Group) {
	for _, g := range groups {
		gg := g
		newGroups = append(newGroups, &gg)
	}
	return
}

type Attack struct {
	Attacker *Group
	Defender *Group
}

func (a Attack) String() string {
	return fmt.Sprintf("(%s%d) -> (%s%d) %d", a.Attacker.Type, a.Attacker.ID, a.Defender.Type, a.Defender.ID, a.Attacker.Dammage(*a.Defender))
}

func fight(immune, infections []*Group, boost int) (int, GroupType) {
	for _, g := range immune {
		g.AttackDammage += boost
	}
	for len(immune) > 0 && len(infections) > 0 {
		turnOrder := make([]*Group, len(immune), len(immune)+len(infections))
		copy(turnOrder, immune)
		turnOrder = append(turnOrder, infections...)
		sortByEffectivePower(turnOrder)

		attackedBy := make(map[*Group]*Group)
		for _, attackingGroup := range turnOrder {
			// selection targets
			var defenders []*Group
			if attackingGroup.Type == Immune {
				defenders = infections
			} else {
				defenders = immune
			}
			var maxGroup *Group
			var maxDammage int
			for _, defendingGroup := range defenders {
				if _, ok := attackedBy[defendingGroup]; ok {
					// already chosen as target
					continue
				}
				potentialAttack := attackingGroup.Dammage(*defendingGroup)
				if potentialAttack != 0 {
					if potentialAttack == maxDammage {
						if defendingGroup.EP() == maxGroup.EP() {
							if defendingGroup.Initiative > maxGroup.Initiative {
								maxGroup = defendingGroup
							}
						} else if defendingGroup.EP() > maxGroup.EP() {
							maxGroup = defendingGroup
						}
					} else if potentialAttack > maxDammage {
						maxDammage = potentialAttack
						maxGroup = defendingGroup
					}
				}
			}
			if maxGroup != nil {
				attackedBy[maxGroup] = attackingGroup
			}
		}
		attacks := make([]Attack, 0, len(attackedBy))
		for defender, attacker := range attackedBy {
			attacks = append(attacks, Attack{Attacker: attacker, Defender: defender})
		}
		sortByInitiative(attacks)
		var maxKilled int
		for _, attack := range attacks {
			killedUnits := attack.Attacker.Dammage(*attack.Defender) / attack.Defender.HP
			if killedUnits > attack.Defender.Units {
				killedUnits = attack.Defender.Units
			}
			if killedUnits > maxKilled {
				maxKilled = killedUnits
			}
			attack.Defender.Units -= killedUnits
		}
		immune = removeKilled(immune)
		infections = removeKilled(infections)
		if maxKilled == 0 {
			return 0, Infection
		}
	}
	var result int
	var winner GroupType
	for _, g := range immune {
		result += g.Units
		winner = g.Type
	}
	for _, g := range infections {
		result += g.Units
		winner = g.Type
	}
	return result, winner
}

func fightUntilVictory(immune, infections []Group) int {
	winner := Infection
	var boost, result int
	for winner == Infection {
		result, winner = fight(copyGroups(immune), copyGroups(infections), boost)
		fmt.Println(boost, result, winner)
		boost++
	}
	return result
}

func sortByEffectivePower(groups []*Group) {
	sort.Slice(groups, func(i, j int) bool {
		if groups[i].EP() == groups[j].EP() {
			return groups[i].Initiative > groups[j].Initiative
		} else {
			return groups[i].EP() > groups[j].EP()
		}
	})
}

func sortByInitiative(attacks []Attack) {
	sort.Slice(attacks, func(i, j int) bool {
		return attacks[i].Attacker.Initiative > attacks[j].Attacker.Initiative
	})
}

func removeKilled(groups []*Group) []*Group {
	newGroups := groups[:0]
	for _, g := range groups {
		if g.Units > 0 {
			newGroups = append(newGroups, g)
		}
	}
	return newGroups
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}
