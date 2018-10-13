#include<iostream>
#include<stdlib.h>

using namespace std;

bool checkdivisibility(int p, int lb, int ub) {
  for (int i = lb; i < ub; i++) {
    if (p % i == 0) {
      return true;
    }
  }
  return false;
}

int main(int argc, char* argv[]) {
  int lb = atoi(argv[1]);
  int ub = atoi(argv[2]);
  if (checkdivisibility(81, lb, ub)) {
    cout << "true" << endl;
  } else {
    cout << "false" << endl;
  }
  return 0;
}
