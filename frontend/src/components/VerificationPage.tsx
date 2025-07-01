import React, { useState, useEffect } from 'react';
import { Container, Box, Typography, TextField, Button, Alert } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const VerificationPage = () => {
  const navigate = useNavigate();
  const [code, setCode] = useState('');
  const [error, setError] = useState('');

  const [focusedIndex, setFocusedIndex] = useState(0);

  useEffect(() => {
    const input = document.getElementById(`code-input-${focusedIndex}`) as HTMLInputElement | null;
    input?.focus();
  }, [focusedIndex]);

  const handleVerify = async () => {
    try {
      const phoneNumber = localStorage.getItem('adminPhoneNumber');
      if (!phoneNumber) {
        setError('Phone number not found');
        return;
      }

      const response = await axios.post('http://localhost:8080/api/verify-code', {
        phoneNumber,
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

  return (
    <Container maxWidth="sm" sx={{ mt: 4 }}>
      <Box sx={{ textAlign: 'center' }}>
        <Typography variant="h5" component="h1" gutterBottom>
          Verify Your Code
        </Typography>
        <Typography variant="body1" sx={{ mb: 4 }}>
          Enter the 6-digit verification code sent to your phone
        </Typography>
        
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
          {[...Array(6)].map((_, index) => (
            <TextField
              key={index}
              id={`code-input-${index}`}
              value={code[index] || ''}
              onChange={(e) => {
                const newCode = e.target.value;
                if (newCode.length === 1 && newCode.match(/[0-9]/)) {
                  setCode(code.substring(0, index) + newCode + code.substring(index + 1));
                  if (index < 5) {
                    setFocusedIndex(index + 1);
                  }
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
