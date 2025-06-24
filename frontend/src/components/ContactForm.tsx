import React, { useState } from 'react';
import {
  Box,
  Button,
  Container,
  TextField,
  Typography,
  Alert,
} from '@mui/material';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import axios from 'axios';

const contactSchema = z.object({
  name: z.string()
    .min(2, 'Name must be at least 2 characters')
    .max(100, 'Name must be at most 100 characters'),
  email: z.string()
    .email('Email must be valid')
    .min(1, 'Email is required'),
  phone: z.string()
    .min(10, 'Phone number must be at least 10 digits')
    .max(20, 'Phone number must be at most 20 digits'),
  message: z.string()
    .max(500, 'Message must be at most 500 characters')
    .optional(),
});

type ContactFormInputs = z.infer<typeof contactSchema>;

const ContactForm = () => {
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState('');

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<ContactFormInputs>({
    resolver: zodResolver(contactSchema),
    defaultValues: {
      name: '',
      email: '',
      phone: '',
      message: '',
    },
  });

  const onSubmit = async (data: ContactFormInputs) => {
    try {
      const response = await axios.post('http://localhost:8080/api/contact', data);
      if (response.status === 200) {
        setSuccess(response.data.message || 'Your message has been sent successfully!');
        reset();
      }
    } catch (err) {
      setError(err.response?.data?.error || 'An error occurred');
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Contact Us
        </Typography>

        {success && (
          <Alert
            onClose={() => setSuccess(false)}
            severity="success" sx={{ mb: 2 }}
          >
            Your message has been sent successfully!
          </Alert>
        )}

        {error && (
          <Alert
            onClose={() => setError('')}
            severity="error" sx={{ mb: 2 }}
          >
            {error}
          </Alert>
        )}

        <form onSubmit={handleSubmit(onSubmit)}>
          <TextField
            fullWidth
            id="name"
            label="Name"
            {...register('name')}
            error={!!errors.name}
            helperText={errors.name?.message}
            sx={{ mb: 2 }}
          />

          <TextField
            fullWidth
            id="email"
            label="Email"
            type="email"
            {...register('email')}
            error={!!errors.email}
            helperText={errors.email?.message}
            sx={{ mb: 2 }}
          />

          <TextField
            fullWidth
            id="phone"
            label="Phone Number"
            {...register('phone')}
            error={!!errors.phone}
            helperText={errors.phone?.message}
            sx={{ mb: 2 }}
          />

          <TextField
            fullWidth
            id="message"
            label="Message"
            multiline
            rows={4}
            {...register('message')}
            error={!!errors.message}
            helperText={errors.message?.message}
            sx={{ mb: 2 }}
          />

          <Button
            type="submit"
            variant="contained"
            color="primary"
            fullWidth
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Sending...' : 'Send Message'}
          </Button>
        </form>
      </Box>
    </Container>
  );
};

export default ContactForm;
