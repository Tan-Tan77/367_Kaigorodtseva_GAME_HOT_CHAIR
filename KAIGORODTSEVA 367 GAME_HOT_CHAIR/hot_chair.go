package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Character interface {
	Hit()
	Block()
}

var parts = map[int]string{
	1: "Туловище",
	2: "Голова",
}

type ItemType int

const (
	Weapon ItemType = iota
	Armor
	Consumable
)

func (it ItemType) String() string {
	switch it {
	case Weapon:
		return "Оружие"
	case Armor:
		return "Броня"
	case Consumable:
		return "Расходник"
	default:
		return "Неизвестно"
	}
}

type Item struct {
	Name    string
	Type    ItemType
	Attack  int
	Defence int
	PlusHP  int
	Slot    string
}

var weapons = []Item{
	{"Ржавый меч", Weapon, 1, 0, 0, "right_hand"},
	{"Острые когти", Weapon, 2, 0, 0, "right_hand"},
	{"Меч защитника", Weapon, 3, 0, 0, "right_hand"},
	{"Пламенный клинок", Weapon, 4, 0, 0, "right_hand"},
}

var armors = []Item{
	{"Железные перчатки", Armor, 0, 1, 0, "hands"},
	{"Пламенная броня", Armor, 0, 3, 0, "chest"},
	{"Щит из драконьей чешуи", Armor, 0, 2, 0, "left_hand"},
}

var consumables = []Item{
	{"Аптечка", Consumable, 0, 0, 2, "consumable"},
	{"Защитное зелье", Consumable, 0, 1, 0, "consumable"},
}

type Player struct {
	Name        string
	HP          int
	Streight    int
	hit         string
	block       string
	Inventory   []Item
	Equipment   map[string]Item
	BaseAttack  int
	BaseDefence int
}

type Enemy struct {
	Name     string
	HP       int
	Streight int
	hit      string
	block    string
	Item     *Item
}

func (p *Player) HitEnemyPart(e *Enemy, part string) {
	e.hit = part
	fmt.Printf("Игрок %s выбрал бить в часть тела %s \n", p.Name, e.hit)
}

func (e *Enemy) HitPlayerPart(p *Player, part string) {
	p.hit = part
	fmt.Printf("Враг %s выбрал бить в часть тела %s \n", e.Name, p.hit)
}

func (e *Enemy) BlockEnemy(i int) {
	if i != 0 {
		e.block = parts[i]
		fmt.Printf("Враг %s выбрал защищать часть тела %s \n", e.Name, e.block)
	} else if i == 0 {
		fmt.Printf("Враг %s не защитился \n", e.Name)
	}
}

func (p *Player) BlockPlayer(i int) {
	if i != 0 {
		p.block = parts[i]
		fmt.Printf("Игрок %s выбрал защищать часть тела %s \n", p.Name, p.block)
	} else if i == 0 {
		fmt.Printf("Игрок %s не защитился \n", p.Name)
	}
}

func (p *Player) GetCurrentAttack() int {
	total := p.BaseAttack
	for _, item := range p.Equipment {
		total += item.Attack
	}
	return total
}

func (p *Player) GetCurrentDefence() int {
	total := p.BaseDefence
	for _, item := range p.Equipment {
		total += item.Defence
	}
	return total
}

func (p *Player) Hit(e *Enemy) {
	var amount int
	currentAttack := p.GetCurrentAttack()

	if e.block != "" {
		if e.HP > 0 {
			if e.block == parts[1] {
				damage := currentAttack - 1
				if damage < 0 {
					damage = 0
				}
				amount = e.HP - damage
				e.HP = amount
				fmt.Printf("Враг %v получил на 1 единицу урона меньше\n", e.Name)
				fmt.Printf("Игрок %v нанёс урона %d врагу %v \n", p.Name, damage, e.Name)
			} else if e.block == parts[2] {
				damage := currentAttack - 2
				if damage < 0 {
					damage = 0
				}
				amount = e.HP - damage
				e.HP = amount
				fmt.Printf("Враг %v получил на 2 единицы урона меньше\n", e.Name)
				fmt.Printf("Игрок %v нанёс урона %d врагу %v \n", p.Name, damage, e.Name)
			}
		}
	} else {
		amount = e.HP - currentAttack
		e.HP = amount
		fmt.Printf("Игрок %v нанёс урона %d врагу %v \n", p.Name, currentAttack, e.Name)
	}

	if e.HP <= 0 {
		fmt.Printf("Враг мёртв %s\n", e.Name)

		if e.Item != nil {
			fmt.Printf("Вы получили трофей: %s!\n", e.Item.Name)
			p.AddToInventory(*e.Item)
			e.Item = nil
		}
	}
}

func (p *Player) HitPlayer(target *Player) {
	var amount int
	currentAttack := p.GetCurrentAttack()

	if target.block != "" {
		if target.HP > 0 {
			if target.block == parts[1] {
				damage := currentAttack - 1
				if damage < 0 {
					damage = 0
				}
				amount = target.HP - damage
				target.HP = amount
				fmt.Printf("Игрок %v получил на 1 единицу урона меньше (защита брони: %d)\n", target.Name, target.GetCurrentDefence())
				fmt.Printf("Игрок %v нанёс урона %d игроку %v \n", p.Name, damage, target.Name)
			} else if target.block == parts[2] {
				damage := currentAttack - 2
				if damage < 0 {
					damage = 0
				}
				amount = target.HP - damage
				target.HP = amount
				fmt.Printf("Игрок %v получил на 2 единицы урона меньше (защита брони: %d)\n", target.Name, target.GetCurrentDefence())
				fmt.Printf("Игрок %v нанёс урона %d игроку %v \n", p.Name, damage, target.Name)
			}
		}
	} else {
		damage := currentAttack - target.GetCurrentDefence()
		if damage < 0 {
			damage = 0
		}
		amount = target.HP - damage
		target.HP = amount
		fmt.Printf("Игрок %v нанёс урона %d игроку %v (защита брони: %d)\n", p.Name, damage, target.Name, target.GetCurrentDefence())
	}

	if target.HP <= 0 {
		fmt.Printf("Игрок %s побеждён!\n", target.Name)
	}
}

func (e *Enemy) Hit(p *Player) {
	var amount int
	currentDefence := p.GetCurrentDefence()
	effectiveStrength := e.Streight - currentDefence
	if effectiveStrength < 0 {
		effectiveStrength = 0
	}

	if p.block != "" {
		if p.HP > 0 {
			if p.block == parts[1] {
				damage := effectiveStrength - 1
				if damage < 0 {
					damage = 0
				}
				amount = p.HP - damage
				p.HP = amount
				fmt.Printf("Игрок %v получил на 1 единицу урона меньше (защита брони: %d)\n", p.Name, currentDefence)
				fmt.Printf("Враг %v нанёс урона %d игроку %v \n", e.Name, damage, p.Name)
			} else if p.block == parts[2] {
				damage := effectiveStrength - 2
				if damage < 0 {
					damage = 0
				}
				amount = p.HP - damage
				p.HP = amount
				fmt.Printf("Игрок %v получил на 2 единицы урона меньше (защита брони: %d)\n", p.Name, currentDefence)
				fmt.Printf("Враг %v нанёс урона %d игроку %v \n", e.Name, damage, p.Name)
			}
		}
	} else {
		amount = p.HP - effectiveStrength
		p.HP = amount
		fmt.Printf("Враг %v нанёс урона %d игроку %v (защита брони: %d)\n", e.Name, effectiveStrength, p.Name, currentDefence)
	}

	if p.HP <= 0 {
		fmt.Printf("Игрок мёртв %s\n", p.Name)
	}
}

func (p *Player) AddToInventory(item Item) {
	p.Inventory = append(p.Inventory, item)
	fmt.Printf("Предмет %s добавлен в инвентарь\n", item.Name)
}

func (p *Player) ShowInventory() {
	fmt.Println("\n=== ИНВЕНТАРЬ ===")
	if len(p.Inventory) == 0 {
		fmt.Println("Инвентарь пуст")
		return
	}
	for i, item := range p.Inventory {
		fmt.Printf("%d. %s (%s)", i+1, item.Name, item.Type)
		switch item.Type {
		case Weapon:
			fmt.Printf(" +%d к атаке", item.Attack)
		case Armor:
			fmt.Printf(" +%d к защите", item.Defence)
		case Consumable:
			fmt.Printf(" +%d к здоровью", item.PlusHP)
		}
		fmt.Println()
	}
}

func (p *Player) ShowEquipment() {
	fmt.Println("\n=== ЭКИПИРОВКА ===")
	if len(p.Equipment) == 0 {
		fmt.Println("Нет экипированных предметов")
		return
	}
	for slot, item := range p.Equipment {
		fmt.Printf("%s: %s", slot, item.Name)
		switch item.Type {
		case Weapon:
			fmt.Printf(" (+%d к атаке)", item.Attack)
		case Armor:
			fmt.Printf(" (+%d к защите)", item.Defence)
		}
		fmt.Println()
	}
	fmt.Printf("Итоговая атака: %d\n", p.GetCurrentAttack())
	fmt.Printf("Итоговая защита: %d\n", p.GetCurrentDefence())
}

func (p *Player) TakeOff() {
	if len(p.Equipment) == 0 {
		fmt.Println("Нет экипированных предметов")
		return
	}

	fmt.Println("\nВыберите предмет для снятия:")
	var items []Item
	var slots []string
	i := 1
	for slot, item := range p.Equipment {
		fmt.Printf("%d. %s (%s) - слот: %s\n", i, item.Name, item.Type, slot)
		items = append(items, item)
		slots = append(slots, slot)
		i++
	}
	fmt.Printf("%d. Отмена\n", i)

	var choice int
	fmt.Scan(&choice)

	if choice == i {
		fmt.Println("Отмена")
		return
	}

	if choice < 1 || choice > len(items) {
		fmt.Println("Неверный выбор")
		return
	}

	item := items[choice-1]
	slot := slots[choice-1]

	delete(p.Equipment, slot)
	p.Inventory = append(p.Inventory, item)
	fmt.Printf("Предмет %s снят и помещён в инвентарь\n", item.Name)
}

func (p *Player) Equip() {
	if len(p.Inventory) == 0 {
		fmt.Println("Инвентарь пуст")
		return
	}

	fmt.Println("\nВыберите предмет для экипировки:")
	for i, item := range p.Inventory {
		fmt.Printf("%d. %s (%s)", i+1, item.Name, item.Type)
		switch item.Type {
		case Weapon:
			fmt.Printf(" +%d к атаке", item.Attack)
		case Armor:
			fmt.Printf(" +%d к защите", item.Defence)
		case Consumable:
			fmt.Printf(" +%d к здоровью", item.PlusHP)
		}
		fmt.Println()
	}
	fmt.Printf("%d. Отмена\n", len(p.Inventory)+1)

	var choice int
	fmt.Scan(&choice)

	if choice == len(p.Inventory)+1 {
		fmt.Println("Отмена")
		return
	}

	if choice < 1 || choice > len(p.Inventory) {
		fmt.Println("Неверный выбор")
		return
	}

	item := p.Inventory[choice-1]

	switch item.Type {
	case Weapon, Armor:

		if equipped, exists := p.Equipment[item.Slot]; exists {
			fmt.Printf("Слот %s уже занят предметом: %s\n", item.Slot, equipped.Name)
			fmt.Println("Сначала снимите текущий предмет")
			return
		}

		p.Equipment[item.Slot] = item

		p.Inventory = append(p.Inventory[:choice-1], p.Inventory[choice:]...)
		fmt.Printf("Предмет %s экипирован в слот %s\n", item.Name, item.Slot)

	case Consumable:

		p.HP += item.PlusHP
		fmt.Printf("Использован %s. Здоровье увеличено на %d. Текущее здоровье: %d\n",
			item.Name, item.PlusHP, p.HP)

		p.Inventory = append(p.Inventory[:choice-1], p.Inventory[choice:]...)
	}
}

func generateRandomItem() *Item {
	rand.Seed(time.Now().UnixNano())
	itemType := rand.Intn(3)

	switch itemType {
	case 0:
		weapon := weapons[rand.Intn(len(weapons))]
		return &weapon
	case 1:
		armor := armors[rand.Intn(len(armors))]
		return &armor
	case 2:
		consumable := consumables[rand.Intn(len(consumables))]
		return &consumable
	default:
		return nil
	}
}

func hotSeatMode() {
	fmt.Println("\n=== РЕЖИМ ГОРЯЧЕГО СТУЛА (PvP) ===")
	fmt.Println("Два игрока будут сражаться друг против друга на одном компьютере!")

	players := make([]Player, 2)

	players[0] = Player{
		Name:        "Игрок 1",
		HP:          10,
		Streight:    3,
		BaseAttack:  3,
		BaseDefence: 0,
		Inventory: []Item{
			{"Меч защитника", Weapon, 3, 0, 0, "right_hand"},
			{"Пламенная броня", Armor, 0, 3, 0, "chest"},
			{"Аптечка", Consumable, 0, 0, 2, "consumable"},
		},
		Equipment: make(map[string]Item),
	}

	players[1] = Player{
		Name:        "Игрок 2",
		HP:          10,
		Streight:    3,
		BaseAttack:  3,
		BaseDefence: 0,
		Inventory: []Item{
			{"Острые когти", Weapon, 2, 0, 0, "right_hand"},
			{"Щит из драконьей чешуи", Armor, 0, 2, 0, "left_hand"},
			{"Защитное зелье", Consumable, 0, 1, 0, "consumable"},
		},
		Equipment: make(map[string]Item),
	}

	players[0].Equipment["right_hand"] = players[0].Inventory[0]
	players[0].Equipment["chest"] = players[0].Inventory[1]
	players[0].Inventory = players[0].Inventory[2:]

	players[1].Equipment["right_hand"] = players[1].Inventory[0]
	players[1].Equipment["left_hand"] = players[1].Inventory[1]
	players[1].Inventory = players[1].Inventory[2:]

	fmt.Println("\n=== НАЧАЛЬНЫЕ ХАРАКТЕРИСТИКИ ===")
	for i := 0; i < 2; i++ {
		fmt.Printf("\n%s:\n", players[i].Name)
		fmt.Printf("Здоровье: %d\n", players[i].HP)
		fmt.Printf("Атака: %d\n", players[i].GetCurrentAttack())
		fmt.Printf("Защита: %d\n", players[i].GetCurrentDefence())
	}

	fmt.Println("\n=== ПОДГОТОВКА К БОЮ ===")
	for i := 0; i < 2; i++ {
		fmt.Printf("\n--- Ход игрока %d ---\n", i+1)
		fmt.Println("Вы можете:")
		fmt.Println("1. Показать инвентарь")
		fmt.Println("2. Показать экипировку")
		fmt.Println("3. Надеть предмет")
		fmt.Println("4. Снять предмет")
		fmt.Println("5. Использовать расходник")
		fmt.Println("6. Продолжить")

		var choice int
		for {
			fmt.Print("Ваш выбор: ")
			fmt.Scan(&choice)

			switch choice {
			case 1:
				players[i].ShowInventory()
			case 2:
				players[i].ShowEquipment()
			case 3:
				players[i].Equip()
			case 4:
				players[i].TakeOff()
			case 5:
				if len(players[i].Inventory) > 0 {
					players[i].ShowInventory()
					fmt.Print("Выберите номер расходника для использования (0 для отмены): ")
					var itemChoice int
					fmt.Scan(&itemChoice)
					if itemChoice > 0 && itemChoice <= len(players[i].Inventory) {
						item := players[i].Inventory[itemChoice-1]
						if item.Type == Consumable {
							players[i].HP += item.PlusHP
							fmt.Printf("Использован %s. Здоровье увеличено на %d. Текущее здоровье: %d\n",
								item.Name, item.PlusHP, players[i].HP)
							players[i].Inventory = append(players[i].Inventory[:itemChoice-1], players[i].Inventory[itemChoice:]...)
						} else {
							fmt.Println("Это не расходник!")
						}
					}
				} else {
					fmt.Println("Нет расходников в инвентаре")
				}
			case 6:
				break
			default:
				fmt.Println("Неверный выбор")
			}
			if choice == 6 {
				break
			}

			fmt.Println("\nВы можете:")
			fmt.Println("1. Показать инвентарь")
			fmt.Println("2. Показать экипировку")
			fmt.Println("3. Надеть предмет")
			fmt.Println("4. Снять предмет")
			fmt.Println("5. Использовать расходник")
			fmt.Println("6. Продолжить")
		}
	}

	fmt.Println("\n=== НАЧАЛО БОЯ! ===")

	currentPlayer := 0
	round := 1

	for {
		fmt.Printf("\n=== РАУНД %d ===\n", round)

		fmt.Printf("\nХод игрока: %s\n", players[currentPlayer].Name)
		fmt.Printf("Ваше здоровье: %d\n", players[currentPlayer].HP)
		fmt.Printf("Здоровье противника: %d\n", players[1-currentPlayer].HP)

		fmt.Println("\nВыберите действие:")
		fmt.Println("1. Атаковать")
		fmt.Println("2. Заблокировать")

		var action int
		for {
			fmt.Print("Ваш выбор: ")
			fmt.Scan(&action)

			if action == 1 || action == 2 {
				break
			}
			fmt.Println("Неверный выбор. Введите 1 или 2")
		}

		if action == 1 {
			fmt.Println("\nВыберите часть тела для атаки:")
			fmt.Println("1. Туловище")
			fmt.Println("2. Голова")

			var partChoice int
			for {
				fmt.Print("Ваш выбор: ")
				fmt.Scan(&partChoice)

				if partChoice == 1 || partChoice == 2 {
					break
				}
				fmt.Println("Неверный выбор. Введите 1 или 2")
			}

			players[1-currentPlayer].block = ""
			if partChoice == 1 {
				players[1-currentPlayer].hit = "Туловище"
			} else {
				players[1-currentPlayer].hit = "Голова"
			}

			fmt.Printf("\nХод игрока: %s\n", players[1-currentPlayer].Name)
			fmt.Println("Выберите часть тела для защиты:")
			fmt.Println("0. Не защищаться")
			fmt.Println("1. Туловище")
			fmt.Println("2. Голова")

			var blockChoice int
			for {
				fmt.Print("Ваш выбор: ")
				fmt.Scan(&blockChoice)

				if blockChoice >= 0 && blockChoice <= 2 {
					break
				}
				fmt.Println("Неверный выбор. Введите 0, 1 или 2")
			}

			if blockChoice != 0 {
				players[1-currentPlayer].block = parts[blockChoice]
				fmt.Printf("Игрок %s защищает %s\n", players[1-currentPlayer].Name, players[1-currentPlayer].block)
			} else {
				fmt.Printf("Игрок %s не защищается\n", players[1-currentPlayer].Name)
			}

			players[currentPlayer].HitPlayer(&players[1-currentPlayer])

		} else {
			fmt.Println("\nВыберите часть тела для защиты:")
			fmt.Println("1. Туловище")
			fmt.Println("2. Голова")

			var blockChoice int
			for {
				fmt.Print("Ваш выбор: ")
				fmt.Scan(&blockChoice)

				if blockChoice == 1 || blockChoice == 2 {
					break
				}
				fmt.Println("Неверный выбор. Введите 1 или 2")
			}

			players[currentPlayer].block = parts[blockChoice]
			fmt.Printf("Игрок %s защищает %s\n", players[currentPlayer].Name, players[currentPlayer].block)

			fmt.Printf("\nХод игрока: %s\n", players[1-currentPlayer].Name)
			fmt.Println("Выберите часть тела для атаки:")
			fmt.Println("1. Туловище")
			fmt.Println("2. Голова")

			var partChoice int
			for {
				fmt.Print("Ваш выбор: ")
				fmt.Scan(&partChoice)

				if partChoice == 1 || partChoice == 2 {
					break
				}
				fmt.Println("Неверный выбор. Введите 1 или 2")
			}

			if partChoice == 1 {
				players[currentPlayer].hit = "Туловище"
			} else {
				players[currentPlayer].hit = "Голова"
			}

			players[1-currentPlayer].HitPlayer(&players[currentPlayer])
		}

		if players[0].HP <= 0 || players[1].HP <= 0 {
			break
		}

		players[0].block = ""
		players[1].block = ""

		currentPlayer = 1 - currentPlayer
		round++
	}

	fmt.Println("\n=== БОЙ ОКОНЧЕН! ===")
	if players[0].HP <= 0 && players[1].HP <= 0 {
		fmt.Println("НИЧЬЯ! Оба игрока пали в бою!")
	} else if players[0].HP <= 0 {
		fmt.Printf("ПОБЕДИТЕЛЬ: %s!\n", players[1].Name)
		fmt.Printf("%s побеждает с %d очками здоровья!\n", players[1].Name, players[1].HP)
	} else {
		fmt.Printf("ПОБЕДИТЕЛЬ: %s!\n", players[0].Name)
		fmt.Printf("%s побеждает с %d очками здоровья!\n", players[0].Name, players[0].HP)
	}

	fmt.Println("\nХотите сыграть еще раз в режиме горячего стула? (да/нет)")
	var replay string
	fmt.Scan(&replay)
	if replay == "да" {
		hotSeatMode()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var c Player
	c.Name = "Tsumikasa"
	c.Streight = 4
	c.HP = 7
	c.block = ""
	c.Inventory = []Item{
		{"Ржавый меч", Weapon, 1, 0, 0, "right_hand"},
		{"Кожаный доспех", Armor, 0, 1, 0, "chest"},
		{"Малая аптечка", Consumable, 0, 0, 2, "consumable"},
	}
	c.Equipment = make(map[string]Item)
	c.BaseAttack = c.Streight
	c.BaseDefence = 0

	c.Equipment["right_hand"] = c.Inventory[0]
	c.Equipment["chest"] = c.Inventory[1]
	c.Inventory = c.Inventory[2:]

	var e Enemy
	e.Name = "Evil"
	e.Streight = 2
	e.HP = 4
	e.block = ""
	e.Item = generateRandomItem()

	var e1 Enemy
	e1.Name = "Small Evil 1"
	e1.Streight = 1
	e1.HP = 1
	e1.block = ""
	e1.Item = generateRandomItem()

	var e2 Enemy
	e2.Name = "Small Evil 2 "
	e2.Streight = 1
	e2.HP = 1
	e2.block = ""
	e2.Item = generateRandomItem()

	var e3 Enemy
	e3.Name = "Small Evil 3"
	e3.Streight = 1
	e3.HP = 1
	e3.block = ""
	e3.Item = generateRandomItem()

	fmt.Printf("Каждый день в деревне ничем не отличается от другого. %v - скромный блюститель порядка,глава защитного отряда.\nНо в мирные времена, такое подразделение дремлет...\nСможет ли герой деревни помочь жителям в случае реальной опасности?\n", c.Name)
	fmt.Printf("3 часа 54 ночи. Обычная спокойная ночь, как и всегда. Вдруг крики жителей разрывают привычную тишину.\nДеревня озарилась светом от пожара. Паника охватила большинство людей\n%v вскочил с места ночного дежурства и побежал на крики. Несколько домов было охвачено странным красно-чёрным пламенем.\nПроисхождение этого пожара - явно не случайность и имеет причастность к чему-то сверхъестественному. %v заметил странное тёмное существо, резко выделяющиеся позади одного из горящих домов\n%v мгновенно освобождает свой меч из ножен и бежит за монстром. Однако, тот быстро скрылся во тьме.\n", c.Name, c.Name, c.Name)

	fmt.Printf("С того момента жители деревни не знали покоя. Дома вспыхивали один за одним.\nТаинственные существа угрожали спалить всю деревню до тла, но до сих пор никто не мог сказать кто они или что они.\n%v выслеживал мостров несколько дней и наконец он понял, что их логово - зловещая пещера вдалеке от деревни. Чтож, пора выполнить свой долг защитника деревни.\n", c.Name)
	fmt.Printf("Теперь судьба деревни лежит в руках её защитника... Справится ли %v ?\n", c.Name)

	c.ShowEquipment()
	fmt.Println("--------------------------------")

	fmt.Println("Фаза Первая: Встреча с первым монстром ")
	fmt.Println("Вход в пещеру зловещно намекает на ужасы, поджидающие впереди. Игрок проходит всё глубже и глубже внутрь...")
	fmt.Println("--------------------------------")
	fmt.Printf("Статы монстра: Имя: %v Сила удара: %d Здоровье: %d  \n", e.Name, e.Streight, e.HP)
	fmt.Printf("При блокировании головы игрок/враг получает на 2 единицы урона меньше.\nПри блокировке туловища - на 1 единицу меньше.\nЕсли игрок/враг ничего не заблокировал, то получает полный урон.\n")
	fmt.Println("--------------------------------")

	fmt.Println("Хотите проверить инвентарь перед боем? (да/нет)")
	var response string
	fmt.Scan(&response)
	if response == "да" {
		for {
			fmt.Println("\n1. Показать инвентарь")
			fmt.Println("2. Показать экипировку")
			fmt.Println("3. Надеть предмет")
			fmt.Println("4. Снять предмет")
			fmt.Println("5. Начать бой")
			var choice int
			fmt.Scan(&choice)

			switch choice {
			case 1:
				c.ShowInventory()
			case 2:
				c.ShowEquipment()
			case 3:
				c.Equip()
			case 4:
				c.TakeOff()
			case 5:
				break
			default:
				fmt.Println("Неверный выбор")
			}
			if choice == 5 {
				break
			}
		}
	}

	for i := 0; ; i++ {
		e.Hit(&c)
		if e.HP <= 0 || c.HP <= 0 {
			break
		}

		rand := rand.Intn(3)
		e.BlockEnemy(rand)

		if e.HP <= 0 || c.HP <= 0 {
			break
		}

		c.Hit(&e)
		if e.HP <= 0 || c.HP <= 0 {
			break
		}

		c.BlockPlayer(rand)
		if e.HP <= 0 || c.HP <= 0 {
			break
		}
	}

	fmt.Println("--------------------------------")
	fmt.Printf("Здоровье игрока %v: %d \n", c.Name, c.HP)
	fmt.Println("Монстр 1 повержен. Игрок идёт глубже в подземелье.... ")
	fmt.Println("Фаза Вторая: Встреча с 3 монстрами ")
	fmt.Println("Внезапная атака 3 мостров! ")
	fmt.Printf("Статы монстров: Имена: %v, %v,%v. Сила удара: %d. Здоровье: %d.  \n", e1.Name, e2.Name, e3.Name, e1.Streight, e1.HP)
	fmt.Println("--------------------------------")

	if c.HP < 5 {
		fmt.Println("Хотите использовать предмет из инвентаря? (да/нет)")
		var useItem string
		fmt.Scan(&useItem)
		if useItem == "да" {
			c.ShowInventory()
			fmt.Println("Выберите номер предмета для использования (или 0 для отмены):")
			var itemChoice int
			fmt.Scan(&itemChoice)
			if itemChoice > 0 && itemChoice <= len(c.Inventory) {
				item := c.Inventory[itemChoice-1]
				if item.Type == Consumable {
					c.HP += item.PlusHP
					fmt.Printf("Использован %s. Здоровье увеличено на %d. Текущее здоровье: %d\n",
						item.Name, item.PlusHP, c.HP)
					c.Inventory = append(c.Inventory[:itemChoice-1], c.Inventory[itemChoice:]...)
				}
			}
		}
	}

	for i := 0; ; i++ {
		e1.Hit(&c)
		if c.HP <= 0 {
			break
		}
		e2.Hit(&c)
		if c.HP <= 0 {
			break
		}

		rand := rand.Intn(3)
		c.BlockPlayer(rand)
		if c.HP <= 0 {
			break
		}

		e3.Hit(&c)
		if c.HP <= 0 {
			break
		}

		e1.BlockEnemy(rand)
		if c.HP <= 0 {
			break
		}

		c.Hit(&e1)
		c.Hit(&e2)
		c.Hit(&e3)
		if e3.HP <= 0 && e1.HP <= 0 && e2.HP <= 0 {
			break
		}
	}

	if c.HP <= 4 {
		fmt.Println("---------------------------------------------------------------")
		fmt.Printf("Игрок %v выбился из сил в процессе боя... Необходимо найти лечебный корень для повышения здоровья. Уровень здоровья: %d \n", c.Name, c.HP)
		fmt.Println("---------------------------------------------------------------")

		for i := 0; ; i++ {
			var rand = rand.Intn(5)

			if rand == 3 {
				fmt.Printf("Не удалось найти корень... Герой %v смотрит под другим камнем... \n", c.Name)
			} else if rand == 0 {
				fmt.Printf("Герой %v потярянно оглядывает пещеру.. \n", c.Name)
			} else if rand == 1 {
				fmt.Printf("Герой %v смотрит около лужи... \n", c.Name)
			} else if rand == 2 {
				fmt.Printf("Герой %v ничего не сожет найти... \n", c.Name)
			} else if rand == 4 {
				fmt.Printf("Герой %v наконец-то находит корень! \n", c.Name)
				var add int = 2
				c.HP += add
				fmt.Println("---------------------------------------------------------------")
				fmt.Printf("Здоровье повышено на 2 единицы. Здоровье Игрока  %v:  %d \n", c.Name, c.HP)
				break
			}
		}
		fmt.Println("---------------------------------------------------------------")
	} else {
		fmt.Println("---------------------------------------------------------------")
		fmt.Printf("Здоровье Игрока  %v:  %d \n", c.Name, c.HP)
		fmt.Println("---------------------------------------------------------------")
	}

	var boss Enemy
	boss.block = "2"
	boss.HP = 12
	boss.Name = "Origin of evil"
	boss.Streight = 3
	boss.Item = generateRandomItem()

	fmt.Printf("Конец пути всё ближе и ближе... Игрок %v чувствует запах дыма... Кажется, это конец  \n", c.Name)
	fmt.Println("Фаза Третья: Финальная битва ")
	fmt.Println("Леденящий душу силуэт исполинской тени угрожающе навис...\nРуки потрясывает от адреналина и страха перед источником Зла...")
	fmt.Printf("Статы босса: Имя: %v. Сила удара: %d. Здоровье: %d.  \n", boss.Name, boss.Streight, boss.HP)
	fmt.Printf("Босс всегда защищает голову и игрок будет наносить ему на 2 урона меньше\n")
	fmt.Println("---------------------------------------------------------------")
	fmt.Printf("Внезапно голос,трясущий стены пещеры и звучащий из неотуда говорит:\nЭто конец. Ты ничего не сделаешь. Когда твои крылья надежды оборвутся после проигрыша, ты останешься здесь навсегда и твои геройские потуги будут бесполезны\nВсе жители деревни поплнят ряды моих преспешников.\n")
	fmt.Println("Пещера пошатнулась - босс встал в позу для атаки. Его тело охватывает угрожающее пламя.")
	fmt.Println("---------------------------------------------------------------")

	fmt.Println("Последний шанс проверить инвентарь перед финальной битвой! (да/нет)")
	var finalPrep string
	fmt.Scan(&finalPrep)
	if finalPrep == "да" {
		c.ShowInventory()
		c.ShowEquipment()
		fmt.Println("Хотите что-то изменить? (да/нет)")
		var change string
		fmt.Scan(&change)
		if change == "да" {
			for {
				fmt.Println("1. Надеть предмет")
				fmt.Println("2. Снять предмет")
				fmt.Println("3. Начать бой")
				var prepChoice int
				fmt.Scan(&prepChoice)

				if prepChoice == 1 {
					c.Equip()
				} else if prepChoice == 2 {
					c.TakeOff()
				} else if prepChoice == 3 {
					break
				}
			}
		}
	}

	for i := 0; ; i++ {
		var rand = rand.Intn(3)

		boss.Hit(&c)
		if c.HP <= 0 {
			break
		}

		c.BlockPlayer(rand)
		c.Hit(&boss)
		if boss.HP <= 0 {
			break
		}

		boss.BlockEnemy(2)
		c.Hit(&boss)
		if boss.HP <= 0 {
			break
		}
	}

	if c.HP <= 0 {
		fmt.Println("---------------------------------------------------------------")
		fmt.Printf("Игрок %v испустил свой последний вздох... Холод разливается по телу.... Блеск кристалов сверкнул в глазах последний раз...  \n", c.Name)
		fmt.Printf("Теперь судьба деревни туманна и Зло завербует жителей в свою армию тьмы... \n")
		fmt.Printf("Постепенно дома в деревне сгорели все до тла... Атака монстров была беспощадна - злые существа безжалостно отбирали сущности людей\n")
		fmt.Printf("Очень скоро весь пейзаж деревни состоял из гор пепла, горелых домов и разрухи. Ничего не намекало на признаки хоть одного уцелевшего.\nТем временем, в тёмных переходах пещеры монстров стало куда больше...\n")
		fmt.Println(" ")
		fmt.Printf("                          --~|--=*-Плохая концовка: игрок не смог защитить деревню-*=--|~--\n")
		fmt.Println(" ")
		fmt.Println("                                                  ----THE END NUMBER 1----")
		fmt.Println("  ")
	} else if c.HP > 0 {
		fmt.Println("---------------------------------------------------------------")
		fmt.Printf("Игрок %v  одержал победу... Кажется, теперь, всё должно быть хорошо? Лёгкая улыбка появляется на лице героя. Он действительно смог.\n", c.Name)
		fmt.Printf("В деревне всё стало как прежде. Ничего особенного или страшного. Жители деревни благодарят %v за то, что он всех спас.\n", c.Name)

		fmt.Println("\n=== СОБРАННЫЕ ТРОФЕИ ===")
		if len(c.Inventory) > 0 {
			for _, item := range c.Inventory {
				fmt.Printf("- %s (%s)\n", item.Name, item.Type)
			}
		} else {
			fmt.Println("Нет трофеев")
		}
		fmt.Println(" ")
		fmt.Printf("                            --~|--=*-Хорошая концовка: игрок защитил деревню-*=--|~--\n")
		fmt.Println(" ")
		fmt.Println("                                               ----THE END NUMBER 2----")
		fmt.Println("  ")
		fmt.Println("---------------------------------------------------------------")
	}

	fmt.Println("\nХотите попробовать режим 'Горячий стул' (PvP)? (да/нет)")
	var playHotSeat string
	fmt.Scan(&playHotSeat)

	if playHotSeat == "да" {
		hotSeatMode()
	} else {
		fmt.Println("\nСпасибо за игру!")
	}
}
