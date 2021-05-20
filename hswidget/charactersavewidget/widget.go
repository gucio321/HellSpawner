package charactersavewidget

import (
	"fmt"

	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/gucio321/d2d2s"
	"github.com/gucio321/d2d2s/d2senums"
	"github.com/gucio321/d2d2s/d2sitems"
	"github.com/ianling/giu"
)

const (
	maxCharName = 16
	inputIntW   = 30
)

type widget struct {
	id  string
	d2s *d2d2s.D2S
}

func Create(state []byte, id string, d2s *d2d2s.D2S) giu.Widget {
	return &widget{
		id:  id,
		d2s: d2s,
	}
}

func (p *widget) Build() {
	classList := make([]string, 0)
	for class := d2senums.CharacterClassAmazon; class <= d2senums.CharacterClassAssassin; class++ {
		classList = append(classList, class.String())
	}

	class := int32(p.d2s.Class)

	giu.Layout{
		giu.Label(fmt.Sprintf("Version: %s", p.d2s.Version)),
		giu.Line(
			giu.Label("Name:"),
			giu.InputText("##"+p.id+"CharacterName", &p.d2s.Name).OnChange(func() {
				if len(p.d2s.Name) > maxCharName {
					p.d2s.Name = p.d2s.Name[:maxCharName-1]
				}
			}),
		),
		giu.Combo("##"+p.id+"classList", classList[class], classList, &class).OnChange(func() {
			p.d2s.Class = d2senums.CharacterClass(class)
		}),
		giu.Line(
			giu.Label("Level:"),
			hswidget.MakeInputInt("##"+p.id+"level", inputIntW, &p.d2s.Level, func() {
				p.d2s.SetLevel(p.d2s.Level)
			}),
		),
		giu.Line(
			giu.Label("Time:"),
			hswidget.MakeInputInt("##"+p.id+"time", 80, &p.d2s.Time, nil),
		),
		giu.Label(fmt.Sprintf("Map ID: %v", p.d2s.MapID)),
		giu.TreeNode("Status##" + p.id + "status").Layout(p.makeStatusLayout()),
		giu.TreeNode("Difficulty##" + p.id + "difficultyStatus").Layout(p.makeDifficultyLayout()),
		giu.TreeNode("Mercenary##" + p.id + "merc").Layout(p.makeMercLayout()),
		giu.TreeNode("Quests##" + p.id + "quests").Layout(p.makeQuestsLayout()),
		giu.TreeNode("Waypoints##" + p.id + "wqypoints").Layout(p.makeWaypointsLayout()),
		giu.TreeNode("NPC introductions##" + p.id + "npc").Layout(p.makeNPCLayout()),
		giu.TreeNode("Stats##" + p.id + "stats").Layout(p.makeStatsLayout()),
		giu.TreeNode("Items##" + p.id + "items").Layout(p.makeItemsLayout(p.id+"items", p.d2s.Items)),
	}.Build()
}

func (p *widget) makeStatusLayout() giu.Layout {
	return giu.Layout{
		giu.Checkbox("Hardcore##"+p.id+"isHardcore", &p.d2s.Status.Hardcore),
		giu.Checkbox("Died##"+p.id+"isHardcore", &p.d2s.Status.Died),
		giu.Checkbox("Expansion##"+p.id+"isHardcore", &p.d2s.Status.Expansion),
		giu.Checkbox("Ladder (?)##"+p.id+"isHardcore", &p.d2s.Status.Ladder),
	}
}

func (p *widget) makeDifficultyLayout() giu.Layout {
	state := p.getState()
	act := (*p.d2s.Difficulty)[d2enum.DifficultyType(state.difficultyStatus)].Act
	act++
	return giu.Layout{
		giu.SliderInt("##"+p.id+"difficultylevel", &state.difficultyStatus, 0, 2).Format(fmt.Sprintf("Difficulty level %v", d2enum.DifficultyType(state.difficultyStatus))),
		giu.Checkbox("Active##"+p.id+"difficultyStatusActive", &(*p.d2s.Difficulty)[d2enum.DifficultyType(state.difficultyStatus)].Active),
		hswidget.MakeInputInt("Act##"+p.id+"difficultyStatusAct", inputIntW, &act, func() {
			(*p.d2s.Difficulty)[d2enum.DifficultyType(state.difficultyStatus)].Act = act - 1
		}),
	}
}

func (p *widget) makeMercLayout() giu.Layout {
	return giu.Layout{
		giu.Label(fmt.Sprintf("ID: %v", p.d2s.Mercenary.ID)),
		giu.Label(fmt.Sprintf("Name ID: %v", p.d2s.Mercenary.Name)),
		giu.Label(fmt.Sprintf("Type: %v", p.d2s.Mercenary.Type.Code)),
		giu.Line(
			giu.Label("Experience:"),
			hswidget.MakeInputInt("##"+p.id+"mercExp", 80, &p.d2s.Mercenary.Experience, nil),
		),
		// giu.TreeNode("Items##" + p.id + "MercItems").Layout(p.makeItemsLayout(p.id+"mercItems",p.d2s.Mercenary.Items)),
	}
}

func (p *widget) makeItemsLayout(id string, items *d2sitems.Items) giu.Layout {
	state := p.getState()
	numItems := len(*items)
	item := &((*items)[state.itemIdx])
	return giu.Layout{
		giu.SliderInt("Item##"+id+"itemidx", &state.itemIdx, 0, int32(numItems-1)),
		giu.Label(fmt.Sprintf("Type: %s", item.Type)),
		giu.Label(fmt.Sprintf("\t\t: %s", item.TypeName)),
		giu.Checkbox("Identified##"+id+"isitemidentified", &item.Identified),
	}
}

func (p *widget) makeQuestsLayout() giu.Layout {
	state := p.getState()
	numQuests := len((*p.d2s.Quests)[d2enum.DifficultyType(state.questsDifficulty)][state.questsAct].Quests)
	return giu.Layout{
		giu.SliderInt("##"+p.id+"difficultylevelQuests", &state.questsDifficulty, 0, 2).Format(fmt.Sprintf("Difficulty level %v", d2enum.DifficultyType(state.questsDifficulty))),
		giu.SliderInt("##"+p.id+"actQuests", &state.questsAct, 0, 4).Format(fmt.Sprintf("Act %v", state.questsAct+1)),
		giu.SliderInt("##"+p.id+"idxQuests", &state.questsIdx, 0, int32(numQuests-1)).Format(fmt.Sprintf("Quest %v", state.questsIdx+1)),
		giu.Separator(),
		giu.Checkbox("Completed", &(*p.d2s.Quests)[d2enum.DifficultyType(state.questsDifficulty)][state.questsAct].Quests[state.questsIdx].Completed),
		giu.Checkbox("Done (all requirements completed - need to get reward)", &(*p.d2s.Quests)[d2enum.DifficultyType(state.questsDifficulty)][state.questsAct].Quests[state.questsIdx].Done),
		giu.Checkbox("Started", &(*p.d2s.Quests)[d2enum.DifficultyType(state.questsDifficulty)][state.questsAct].Quests[state.questsIdx].Started),
		giu.Checkbox("Closed (swirling fire animation played)", &(*p.d2s.Quests)[d2enum.DifficultyType(state.questsDifficulty)][state.questsAct].Quests[state.questsIdx].Closed),
		giu.Checkbox("Just completed (in current game)", &(*p.d2s.Quests)[d2enum.DifficultyType(state.questsDifficulty)][state.questsAct].Quests[state.questsIdx].JustCompleted),
	}
}

func (p *widget) makeWaypointsLayout() giu.Layout {
	return giu.Layout{}
}

func (p *widget) makeNPCLayout() giu.Layout {
	return giu.Layout{}
}

func (p *widget) makeStatsLayout() giu.Layout {
	return giu.Layout{
		hswidget.MakeInputInt("Strength##"+p.id+"statsStrength", inputIntW, &p.d2s.Stats.Strength, nil),
		hswidget.MakeInputInt("Energy##"+p.id+"statsEnergy", inputIntW, &p.d2s.Stats.Energy, nil),
		hswidget.MakeInputInt("Dexterity##"+p.id+"statsDexterity", inputIntW, &p.d2s.Stats.Dexterity, nil),
		hswidget.MakeInputInt("Vitality##"+p.id+"statsVitality", inputIntW, &p.d2s.Stats.Vitality, nil),
		giu.Separator(),
		hswidget.MakeInputInt("Unused stats points##"+p.id+"statsStatPoints", inputIntW, &p.d2s.Stats.UnusedStats, nil),
		hswidget.MakeInputInt("Unused skills points##"+p.id+"statsSkillsPoints", inputIntW, &p.d2s.Stats.UnusedSkillPoints, nil),
		giu.Separator(),
		giu.Line(
			giu.Label("Max HP: "),
			giu.InputFloat("##"+p.id+"statsMaxHP", &p.d2s.Stats.MaxHP).Size(60),
			giu.Label(fmt.Sprintf("Current HP: %v", p.d2s.Stats.CurrentHP)),
			giu.Button("Reset##"+p.id+"resetcurrenthp").OnClick(func() {
				p.d2s.Stats.CurrentHP = p.d2s.Stats.MaxHP
			}),
		),
		giu.Line(
			giu.Label("Max Mana: "),
			giu.InputFloat("##"+p.id+"statsMaxMana", &p.d2s.Stats.MaxMana).Size(60),
			giu.Label(fmt.Sprintf("Current Mana: %v", p.d2s.Stats.CurrentMana)),
			giu.Button("Reset##"+p.id+"resetcurrentmana").OnClick(func() {
				p.d2s.Stats.CurrentMana = p.d2s.Stats.MaxMana
			}),
		),
		giu.Line(
			giu.Label("Max Stamina: "),
			giu.InputFloat("##"+p.id+"statsMaxStamina", &p.d2s.Stats.MaxStamina).Size(60),
			giu.Label(fmt.Sprintf("Current Stamina: %v", p.d2s.Stats.CurrentStamina)),
			giu.Button("Reset##"+p.id+"resetcurrentstamina").OnClick(func() {
				p.d2s.Stats.CurrentStamina = p.d2s.Stats.MaxStamina
			}),
		),
		giu.Layout{
			giu.Label(fmt.Sprintf("Experience (?): %v", p.d2s.Stats.Experience)),
			giu.Tooltip("##" + p.id + "whynoeditexperience").Layout(
				giu.Label("We're sorry, experience value is written as uint32, but ImGui supports only int32\n" +
					"because of that, editing experience value is impossible, because it might cause panic"),
			),
		},
		hswidget.MakeInputInt("Gold##"+p.id+"statsGold", 60, &p.d2s.Stats.Gold, nil),
		hswidget.MakeInputInt("Gold in stash##"+p.id+"statsGoldStashed", 60, &p.d2s.Stats.StashedGold, nil),
	}
}
