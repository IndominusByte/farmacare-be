import json, logging, pokepy
from collections import OrderedDict

logging.basicConfig(level=logging.DEBUG)

def storeData():
    client = pokepy.V2Client()
    pokemons = list()

    for poke_index in range(1,21):
        pokedex = OrderedDict()
        try:
            poke = client.get_pokemon(poke_index)[0]
            pokedex = {
                "name": poke.name,
                "height": poke.height,
                "weight": poke.weight,
                "score": 0,
                "total_match": 0,
                "moves": [],
                "stats": [],
                "type": []
            }
            for stats in poke.stats:
                pokedex["stats"].append({
                    "name": stats.stat.name,
                    "value": stats.base_stat,
                })
            for move in poke.moves:
                m = client.get_move(move.move.name)[0]
                pokedex["moves"].append({
                    "name": m.name,
                    "pp": m.pp,
                    "power": m.power,
                    "accuracy": m.accuracy,
                    "type": m.type.name
                })
            for types in poke.types:
                t = client.get_type(types.type.name)[0]
                pokedex["type"].append({
                    "name": t.name,
                    "double_damage_from": [x.name for x in t.damage_relations.double_damage_from],
                    "double_damage_to": [x.name for x in t.damage_relations.double_damage_to],
                    "half_damage_from": [x.name for x in t.damage_relations.half_damage_from],
                    "half_damage_to": [x.name for x in t.damage_relations.half_damage_to],
                    "no_damage_from": [x.name for x in t.damage_relations.no_damage_from],
                    "no_damage_to": [x.name for x in t.damage_relations.no_damage_to],
                })
            pokemons.append(pokedex)
        except Exception as e:
            logging.debug("Exception %s" % e)

    with open('pokemon.json', 'w', encoding='utf-8') as f:
        json.dump(pokemons, f, ensure_ascii=False, indent=2)


if __name__ == '__main__':
    storeData()
