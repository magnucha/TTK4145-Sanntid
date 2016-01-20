// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

int i = 0;
pthread_mutex_t counter_lock = PTHREAD_MUTEX_INITIALIZER;

void* function1(){
	int j;
	pthread_mutex_lock(&counter_lock);
	for(j = 0; j < 1000000; j++){
		i++;
	}
	pthread_mutex_unlock(&counter_lock);
	return NULL;
}

void* function2(){
	int k;
	pthread_mutex_lock(&counter_lock);
	for(k = 0; k < 1000000; k++){
		i--;
	}
	pthread_mutex_unlock(&counter_lock);
	return NULL;
}

int main(){
	pthread_t thread_1;
	pthread_t thread_2;

	pthread_create(&thread_1, NULL, function1, NULL);
	pthread_create(&thread_2, NULL, function2, NULL);
	
	pthread_mutex_destroy(&counter_lock);
	pthread_join(thread_1, NULL);
	pthread_join(thread_2, NULL);
	printf("%d\n", i);
}
