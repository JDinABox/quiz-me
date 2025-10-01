type question = {
	question: string;
	answers: { [key: number]: string };
	correct: number[];
	explanation: string;
};

type quizObject = {
	id: string;
	name: string;
	description: string;
	maxQuestionsPerQuiz: number;
	questions: question[];
};
document.addEventListener('alpine:init', () => {
	window.Alpine.data('quiz', (q: quizObject) => ({
		id: q.id,
		name: q.name,
		description: q.description,
		maxQuestionsPerQuiz: q.maxQuestionsPerQuiz,
		questions: q.questions,
		currentQuestion: 0,
		score: 0,
		selectedAnswers: [] as string[],
		active: true,
		toggleAnswer(index: string) {
			if (this.selectedAnswers.includes(index)) {
				const i = this.selectedAnswers.indexOf(index);
				if (i > -1) {
					this.selectedAnswers.splice(i, 1);
				}
				return;
			}
			this.selectedAnswers.push(index);
		},
		get correctAnswers() {
			return this.questions[this.currentQuestion].correct;
		},
		get numberOfCorrect() {
			return this.correctAnswers.length;
		},
		get leftToSelect() {
			return this.numberOfCorrect - this.selectedAnswers.length;
		},
		get incorrect() {
			return this.correctAnswers.filter((x) => !this.selectedAnswers.includes(String(x))).length;
		},
		shuffleArray(array: any[]) {
			for (let i = array.length - 1; i > 0; i--) {
				const j = Math.floor(Math.random() * (i + 1));
				[array[i], array[j]] = [array[j], array[i]];
			}
		},
		buttonColor(index: string) {
			if (!this.selectedAnswers.includes(index)) {
				if (this.leftToSelect <= 0 && this.correctAnswers.includes(Number(index))) {
					return 'green-dark';
				}
				return 'gray';
			}
			if (this.leftToSelect > 0) {
				return 'blue';
			}
			if (this.correctAnswers.includes(Number(index))) {
				return 'green';
			}
			return 'red';
		},
		nextQuestion() {
			if (this.incorrect === 0) {
				this.score += 1;
			} else if (this.numberOfCorrect > 1) {
				this.score += (this.numberOfCorrect - this.incorrect) / this.numberOfCorrect;
			}
			if (this.currentQuestion === this.maxQuestionsPerQuiz - 1) {
				this.active = false;
				return;
			}
			this.selectedAnswers = [];
			this.currentQuestion++;
		},
		resetQuclickHandleriz() {
			this.shuffleArray(this.questions);
			this.score = 0;
			this.selectedAnswers = [];
			this.currentQuestion = 0;
			this.active = true;
		},
		init() {
			this.shuffleArray(this.questions);
		}
	}));
});
