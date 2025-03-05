module github.com/Demilade-Omitiran/pokedex

go 1.22.3

require internal/pokeapi v1.0.0

replace internal/pokeapi => ./internal/pokeapi

require internal/pokecache v1.0.0

replace internal/pokecache => ./internal/pokecache
