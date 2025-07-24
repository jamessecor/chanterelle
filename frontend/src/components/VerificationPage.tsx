import React, { useState, useEffect } from 'react';
import { Container, Box, Typography, TextField, Button, Alert } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const VerificationPage = () => {
  const navigate = useNavigate();
  const [code, setCode] = useState<string>('');
  const [error, setError] = useState('');

  const [focusedIndex, setFocusedIndex] = useState(0);

  useEffect(() => {
    const input = document.getElementById(`code-input-${focusedIndex}`) as HTMLInputElement | null;
    input?.focus();
  }, [focusedIndex]);

  const handleVerify = async () => {
    setError('');

    // Validate code format
    if (!/^[0-9]{6}$/.test(code)) {
      setError('Verification code must be 6 digits');
      return;
    }

    try {
      const email = localStorage.getItem('adminEmail');
      if (!email) {
        setError('Email not found');
        return;
      }

      const response = await axios.post(`${import.meta.env.VITE_API_BASE_ADDRESS}/api/verify-code`, {
        email,
        code: code
      });

      if (response.status === 200) {
        // Store token in localStorage
        localStorage.setItem('token', response.data.token);
        navigate('/admin');
      }
    } catch (error) {
      setError('Invalid verification code');
    }
  };

  const updateCode = (index: number, value: string) => {
    const newCode = code.split('');
    newCode[index] = value;
    setCode(newCode.join(''));
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Button variant="contained" color="primary" onClick={() => navigate('/')}>
          Back to Home
        </Button>
        <Typography component="h1" variant="h5">
          Enter Verification Code
        </Typography>
        <Typography variant="body2" sx={{ color: 'text.secondary', mb: 2 }}>
          If the email was valid, you'll receive a verification code to enter
        </Typography>
        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        <Box sx={{ mt: 4, display: 'flex', gap: 1 }}>
          {[0, 1, 2, 3, 4, 5].map((index) => (
            <TextField
              key={index}
              id={`code-input-${index}`}
              value={code[index] || ''}
              onChange={(e) => {
                updateCode(index, e.target.value);
                if (e.target.value && index < 5) {
                  setFocusedIndex(index + 1);
                }
              }}
              onFocus={() => setFocusedIndex(index)}
              inputProps={{ maxLength: 1 }}
              sx={{ flex: 1 }}
              type="text"
              variant="outlined"
              size="small"
              error={error ? true : false}
            />
          ))}
        </Box>

        <Button
          variant="contained"
          color="primary"
          onClick={handleVerify}
          disabled={code.length !== 6}
          fullWidth
        >
          Verify
        </Button>
      </Box>
    </Container>
  );
};

export default VerificationPage;
