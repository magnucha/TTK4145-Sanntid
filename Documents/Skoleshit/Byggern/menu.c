#include "menu.h"

enum _controller_ {PS3, JOYSTICK};
static _controller controller;

void controller_ps3() {
	controller = PS3;
}

void controller_joystick() {
	controller = JOYSTICK;
}

void user_select(uint8_t print_highscores) {
	return;
}













//--------------------------------------------------------------------------

Menu* main_menu;

struct Menu* Menu_new(char* name, uint8_t num_menus, void (*voidfunc)(void), void (*intfunc)(uint8_t)) { 
  struct Menu* p = malloc(sizeof(struct Menu));
  p->name = name;
  p->num_menus = num_menus;
  p->voidfunc = voidfunc;
  p->intfunc = intfunc;
  if (num_menus) {
  	p->next = malloc(sizeof(Menu*)*num_menus);
  }
  return p;
}

void Menu_initialize_lists(Menu* menu) {
	for (uint8_t i = 0; i<menu->num_menus; i++) {
		menu->menus[i]->prev = menu;
		Menu_initialize_lists(menu);
	}
}

void Menu_initialize() {
	main_menu = Menu_new("Main menu", 4, NULL, NULL);
	main_menu->next[0] = Menu_new("Controller", 3, NULL, NULL);
	main_menu->next[0]->next[0] = Menu_new("PS3", 0, &controller_ps3(), NULL);
	main_menu->next[0]->next[1] = Menu_new("Joystick", 0, &controller_joystick(), NULL);
	main_menu->next[1] = Menu_new("Select user", get_num_users(), NULL, &user_select(0));
	main_menu->next[2] = Menu_new("Play", 0, &game_play(), NULL);
	main_menu->next[3] = Menu_new("Highscores", 0, NULL, &user_select(1));
	menu_initialize_lists(main_menu);
}

void Menu_navigate(Menu* menu) {
	while (1) {
		for (int line=0;line <8;line++) {
			OLED_pos(line,line ? 2 : 0); //Offset alle valgmuligheter
			(cursorPos == line) ? OLED_printfInverted(menu->next[line]->name) : OLED_printf(menu->next[line]->name);
		}
		while (joystick_read(0).y_dir || SW3_pressed());	//Passer på at man bare kan flytte ett hakk av gangen
		while (!(joystick_read(0).y_dir || SW3_pressed())); //Holder deg igjen til du gjør noe
		if ((joystick_read(0).y_dir == DOWN) && (cursorPos < numItems)) {
			cursorPos++;
		}
		else if ((joystick_read(0).y_dir == UP) && (cursorPos > 1)) {
			cursorPos--;
		}
		else if (SW3_pressed(0)) {
			if (cursorPos == 0) {
				Menu_navigate(menu->prev);
			}
			else if (cursorPos < menu->num_menus) {
				Menu_navigate(menu->next[cursorPos]);
			}
			else if (!!(menu->voidfunc)) {//Hvis menyen har en funksjon som ikke tar inn noe
				menu->(*voidfunc)();			
			}
			else if (!!(menu->intfunc)) {//Hvis menyen har en funksjon som tar inn en uint8_t
				menu->(*intfunc)(cursorPos);
			}
			return;
		}
		_delay_ms(100);
	}
}
