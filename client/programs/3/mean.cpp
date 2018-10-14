#include<iostream.h>

using namespace std;

int main(int argc, char* argv) {
  int total_sum = 0;
  for (int i = 0; i < argc; i++) {
    total_sum += argv[i];
  }
  float avg = total_sum / argc;
  cout <<  avg << endl;
  return 0;
}

