import {
  Anchor,
  Button,
  Checkbox,
  Group,
  Paper,
  PasswordInput,
  Stack,
  TextInput,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { upperFirst, useToggle } from '@mantine/hooks';

interface AuthenticationFormProps {
  handleRegister: (
    name: string,
    email: string,
    password: string,
    terms: boolean,
  ) => void;
  handleLogin: (email: string, password: string) => void;
}
export function AuthenticationForm({
  handleRegister,
  handleLogin,
}: AuthenticationFormProps) {
  const [type, toggle] = useToggle(['login', 'register']);
  const form = useForm({
    initialValues: {
      email: 'test@gmail.com',
      name: '',
      password: 'testing',
      terms: true,
    },

    validate: {
      email: val => (/^\S+@\S+$/.test(val) ? null : 'Invalid email'),
      password: val =>
        val.length <= 6
          ? 'Password should include at least 6 characters'
          : null,
    },
  });

  return (
    <Paper radius="md" p="xl">
      <form
        onSubmit={form.onSubmit(() => {
          if (type === 'register') {
            handleRegister(
              form.values.name,
              form.values.email,
              form.values.password,
              form.values.terms,
            );
            return;
          }
          if (type == 'login') {
            handleLogin(form.values.email, form.values.password);
          }
        })}
      >
        <Stack>
          {type === 'register' && (
            <TextInput
              label="Name"
              placeholder="Your name"
              value={form.values.name}
              onChange={event => {
                form.setFieldValue('name', event.currentTarget.value);
              }}
              radius="md"
            />
          )}

          <TextInput
            required
            label="Email"
            placeholder="hello@email.com"
            value={form.values.email}
            onChange={event => {
              form.setFieldValue('email', event.currentTarget.value);
            }}
            error={form.errors.email && 'Invalid email'}
            radius="md"
          />

          <PasswordInput
            required
            label="Password"
            placeholder="Your password"
            value={form.values.password}
            onChange={event => {
              form.setFieldValue('password', event.currentTarget.value);
            }}
            error={
              form.errors.password &&
              'Password should include at least 6 characters'
            }
            radius="md"
          />

          {type === 'register' && (
            <Checkbox
              label="I accept terms and conditions"
              checked={form.values.terms}
              onChange={event => {
                form.setFieldValue('terms', event.currentTarget.checked);
              }}
            />
          )}
        </Stack>

        <Group justify="space-between" mt="xl">
          <Anchor
            component="button"
            type="button"
            c="dimmed"
            onClick={() => {
              toggle();
            }}
            size="xs"
          >
            {type === 'register'
              ? 'Already have an account? Login'
              : "Don't have an account? Register"}
          </Anchor>
          <Button type="submit" radius="xl">
            {upperFirst(type)}
          </Button>
        </Group>
      </form>
    </Paper>
  );
}
export default AuthenticationForm;
