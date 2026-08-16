package mutation

type PreimageState string

const (
	PreimageMatches          PreimageState = "matches"
	PreimageDrifted          PreimageState = "drifted"
	PreimageAlreadySatisfied PreimageState = "already_satisfied"
)

func ClassifyPreimage(staged, live, intended []byte) (PreimageState, error) {
	liveIntended, err := Equivalent(live, intended)
	if err != nil {
		return "", err
	}
	if liveIntended {
		return PreimageAlreadySatisfied, nil
	}
	equal, err := Equivalent(staged, live)
	if err != nil {
		return "", err
	}
	if !equal {
		return PreimageDrifted, nil
	}
	return PreimageMatches, nil
}
