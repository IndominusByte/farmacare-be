package handler_http

import (
	"fmt"
	"math"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
0.5 * (attacker.power * (attacker.accuracy/100)) * attacker.stats.attack /
    defender.stats.defense * (double_damage (25%), default 10%)

- all damage dibagi 2 jika ada half_damage
- zero damage if no_damage in move
*/

type Pokemon struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Height     int                `json:"height" bson:"height"`
	Weight     int                `json:"weight" bson:"weight"`
	Moves      []PokemonMove      `json:"moves" bson:"moves"`
	Stats      []PokemonStat      `json:"stats" bson:"stats"`
	Type       []PokemonType      `json:"type" bson:"type"`
	Score      int                `json:"score" bson:"score"`
	TotalMatch int                `json:"total_match" bson:"total_match"`
	Hp         int                `json:"-"`
}

type PokemonMove struct {
	Name          string `json:"name" bson:"name"`
	Pp            int    `json:"pp" bson:"pp"`
	Power         int    `json:"power" bson:"power"`
	Accuracy      int    `json:"accuracy" bson:"accuracy"`
	Type          string `json:"type" bson:"type"`
	DamageToEnemy int    `json:"-"`
}

type PokemonStat struct {
	Name  string `json:"name" bson:"name"`
	Value int    `json:"value" bson:"value"`
}

type PokemonType struct {
	Name             string   `json:"name" bson:"name"`
	DoubleDamageFrom []string `json:"double_damage_from" bson:"double_damage_from"`
	DoubleDamageTo   []string `json:"double_damage_to" bson:"double_damage_to"`
	HalfDamageFrom   []string `json:"half_damage_from" bson:"half_damage_from"`
	HalfDamageTo     []string `json:"half_damage_to" bson:"half_damage_to"`
	NoDamageFrom     []string `json:"no_damage_from" bson:"no_damage_from"`
	NoDamageTo       []string `json:"no_damage_to" bson:"no_damage_to"`
}

type StabAttacker struct {
	Name           string
	DoubleDamageTo []string
	HalfDamageTo   []string
	NoDamageTo     []string
	IsDoubleDamage bool
}

type PokemonFightScore struct {
	Name  string `json:"name" bson:"name"`
	Score int    `json:"score" bson:"score"`
}

type PokemonFightHistory struct {
	Attacker string                      `json:"attacker" bson:"attacker"`
	Move     []string                    `json:"move" bson:"move"`
	Detail   []PokemonFightHistoryDetail `json:"detail" bson:"detail"`
}

type PokemonFightHistoryDetail struct {
	Name           string `json:"name" bson:"name"`
	HpStart        int    `json:"hp_start" bson:"hp_start"`
	HpEnd          int    `json:"hp_end" bson:"hp_end"`
	AttackRecieved int    `json:"attack_recieved" bson:"attack_recieved"`
}

type PokemonOut struct {
	Id            primitive.ObjectID    `json:"id" bson:"_id"`
	BattleAt      primitive.DateTime    `json:"battle_at" bson:"battle_at"`
	PokemonScore  []PokemonFightScore   `json:"pokemon_score" bson:"pokemon_score"`
	BattleHistory []PokemonFightHistory `json:"battle_history" bson:"battle_history"`
}

func getBaseAttackDefend(pokemon *Pokemon, name string) int {
	for _, val := range pokemon.Stats {
		if val.Name == name {
			return val.Value
		}
	}
	return 0
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func BestMove(attacker, defender Pokemon) PokemonMove {
	baseAttack := getBaseAttackDefend(&attacker, "attack")
	baseDefend := getBaseAttackDefend(&defender, "defense")

	typeDefender := map[string][]string{}
	for index, val := range defender.Type {
		if index == 0 {
			typeDefender["double_damage_from"] = val.DoubleDamageFrom
			typeDefender["half_damage_from"] = val.HalfDamageFrom
			typeDefender["no_damage_from"] = val.NoDamageFrom
		} else {
			typeDefender["double_damage_from"] = append(typeDefender["double_damage_from"], val.DoubleDamageFrom...)
			typeDefender["half_damage_from"] = append(typeDefender["half_damage_from"], val.HalfDamageFrom...)
			typeDefender["no_damage_from"] = append(typeDefender["no_damage_from"], val.NoDamageFrom...)
		}
	}

	typeDefender["double_damage_from"] = removeDuplicateStr(typeDefender["double_damage_from"])
	typeDefender["half_damage_from"] = removeDuplicateStr(typeDefender["half_damage_from"])
	typeDefender["no_damage_from"] = removeDuplicateStr(typeDefender["no_damage_from"])

	bestMove := PokemonMove{}
	for _, move := range attacker.Moves {
		if move.Power > 0 && move.Accuracy > 0 {
			doubleDamage := 0.1
			for _, val := range typeDefender["double_damage_from"] {
				if val == move.Type {
					doubleDamage = 0.25
					break
				}
			}

			move.DamageToEnemy = int(math.Floor(
				0.5 * float64((move.Power*(move.Accuracy/100))*baseAttack/baseDefend) * doubleDamage,
			))

			for _, val := range typeDefender["half_damage_from"] {
				if val == move.Type {
					move.DamageToEnemy /= 2
				}
			}

			for _, val := range typeDefender["no_damage_from"] {
				if val == move.Type {
					move.DamageToEnemy = 0
				}
			}

			if bestMove == (PokemonMove{}) {
				bestMove = move
			} else if move.DamageToEnemy > bestMove.DamageToEnemy {
				bestMove = move
			}
		}
	}

	return bestMove
}

func removeIndexPokemon(s []Pokemon, index int) []Pokemon {
	return append(s[:index], s[index+1:]...)
}

func SimulateBattle(pokemons []Pokemon) PokemonOut {
	// set hp outside struct
	for i := 0; i < len(pokemons); i++ {
		for _, val := range pokemons[i].Stats {
			if val.Name == "hp" {
				pokemons[i].Hp = val.Value
				break
			}
		}
	}

	countPokemonKicked, history, score := 0, []PokemonFightHistory{}, []PokemonFightScore{}
	for countPokemonKicked < 4 {
		// set turn pokemon attacker
		var pokeAttacker Pokemon
		for index, value := range pokemons {
			// pop from stack
			if value.Hp > 0 {
				pokeAttacker = value
				pokemons = removeIndexPokemon(pokemons, index)
				break
			}
		}

		// attack defender
		historyMove, historyDetail := []string{}, []PokemonFightHistoryDetail{}
		for i := 0; i < len(pokemons); i++ {
			poke := &pokemons[i]
			if poke.Hp > 0 {
				tmpHp, move := poke.Hp, BestMove(pokeAttacker, *poke)
				if poke.Hp-move.DamageToEnemy < 1 {
					poke.Hp = 0
					countPokemonKicked += 1 // control while loop
					score = append(score, PokemonFightScore{
						Name:  poke.Name,
						Score: countPokemonKicked,
					})
				} else {
					poke.Hp -= move.DamageToEnemy
				}

				historyMove = append(historyMove, fmt.Sprintf("[%d] [%s] - %s", move.DamageToEnemy, move.Name, poke.Name))
				historyDetail = append(historyDetail, PokemonFightHistoryDetail{
					Name:           poke.Name,
					HpStart:        tmpHp,
					HpEnd:          poke.Hp,
					AttackRecieved: move.DamageToEnemy,
				})
			}
		}

		history = append(history, PokemonFightHistory{
			Attacker: pokeAttacker.Name,
			Move:     historyMove,
			Detail:   historyDetail,
		})
		// push pokemon to stack again
		pokemons = append(pokemons, pokeAttacker)
	}

	// search pokemon last man standing
	for _, poke := range pokemons {
		isExist := false
		for _, sc := range score {
			if sc.Name == poke.Name {
				isExist = true
			}
		}
		if !isExist {
			score = append(score, PokemonFightScore{
				Name:  poke.Name,
				Score: 5,
			})
		}
	}

	return PokemonOut{PokemonScore: score, BattleHistory: history}
}
