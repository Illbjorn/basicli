package basicli

func Run[P *T, T any](v P) error {
  if err := Unmarshal(v); err != nil {
    return err
  } else if err = Dispatch(v); err != nil {
    return err
  } else {
    return nil
  }
}
