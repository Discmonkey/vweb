pub trait Node {
    type Input;
    type Output;

    fn call(&mut self, &input: Self::Input) -> Self::Output;
}