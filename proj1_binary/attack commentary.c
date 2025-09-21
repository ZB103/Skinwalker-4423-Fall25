int main(int argc, char** argv)
{
	for (;;) //This runs the loop for infinity until the program is externally killed
		system(argv[0]);
		//system - invokes an operating system command
		//argv[0] - the name of the program itself
}
//I believe this program will repeatedly open new instances of itself. This would quickly take over the computer's memory and CPU. An attacker could use this program to significnantly slow down the user's computer and make it difficult to close the program.
