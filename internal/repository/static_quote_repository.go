package repository

import (
	"github.com/denis-sukhoverkhov/word-of-wisdom/internal/models"
	"math/rand"
	"time"
)

type StaticQuoteRepository struct {
	quotes []models.Quote
	rng    *rand.Rand
}

func NewStaticQuoteRepository() StaticQuoteRepository {
	return StaticQuoteRepository{
		quotes: []models.Quote{
			{Text: "Life is 10% what happens to us and 90% how we react to it.", Author: "Charles R. Swindoll"},
			{Text: "It takes courage to grow up and become who you really are.", Author: "E.E. Cummings"},
			{Text: "Your self-worth is determined by you. You don't have to depend on someone telling you who you are.", Author: "Beyoncé"},
			{Text: "Nothing is impossible. The word itself says 'I'm possible!'", Author: "Audrey Hepburn"},
			{Text: "Keep your face always toward the sunshine, and shadows will fall behind you.", Author: "Walt Whitman"},
			{Text: "You have brains in your head. You have feet in your shoes. You can steer yourself any direction you choose. You're on your own. And you know what you know. And you are the guy who'll decide where to go.", Author: "Dr. Seuss"},
			{Text: "Attitude is a little thing that makes a big difference.", Author: "Winston Churchill"},
			{Text: "To bring about change, you must not be afraid to take the first step. We will fail when we fail to try.", Author: "Rosa Parks"},
			{Text: "All our dreams can come true, if we have the courage to pursue them.", Author: "Walt Disney"},
			{Text: "Don't sit down and wait for the opportunities to come. Get up and make them.", Author: "Madam C.J. Walker"},
			{Text: "Champions keep playing until they get it right.", Author: "Billie Jean King"},
			{Text: "I am lucky that whatever fear I have inside me, my desire to win is always stronger.", Author: "Serena Williams"},
			{Text: "You are never too old to set another goal or to dream a new dream.", Author: "C.S. Lewis"},
			{Text: "It is during our darkest moments that we must focus to see the light.", Author: "Aristotle"},
			{Text: "Believe you can and you're halfway there.", Author: "Theodore Roosevelt"},
			{Text: "Life shrinks or expands in proportion to one’s courage.", Author: "Anaïs Nin"},
			{Text: "Just don't give up trying to do what you really want to do. Where there is love and inspiration, I don't think you can go wrong.", Author: "Ella Fitzgerald"},
			{Text: "Try to be a rainbow in someone's cloud.", Author: "Maya Angelou"},
			{Text: "If you don't like the road you're walking, start paving another one.", Author: "Dolly Parton"},
			{Text: "Real change, enduring change, happens one step at a time.", Author: "Ruth Bader Ginsburg"},
			{Text: "All dreams are within reach. All you have to do is keep moving towards them.", Author: "Viola Davis"},
		},
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r StaticQuoteRepository) GetRandomQuote() (models.Quote, error) {
	return r.quotes[r.rng.Intn(len(r.quotes))], nil
}
