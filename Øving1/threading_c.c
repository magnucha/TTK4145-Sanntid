// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

int i = 0;

void* function1(){
	for(int j = 0; j < 1000000; j++){
		i++;
	}
	return NULL;
}

void* function2(){
	for(int k = 0; k < 1000000; k++){
		i--;
	}
	return NULL;
}

int main(){
	pthread_t thread_1;
	pthread_t thread_2;
	pthread_create(&thread_1, NULL, function1, NULL);
	pthread_create(&thread_2, NULL, function2, NULL);
	pthread_join(thread_1, NULL);
	pthread_join(thread_2, NULL);
	printf("%d\n", i);
}
