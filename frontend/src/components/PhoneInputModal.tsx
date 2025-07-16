import React, { useState } from 'react';
import { Modal, Box, TextField, Button, Typography, Alert } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

interface PhoneInputModalProps {
  open: boolean;
  onClose: () => void;
}

const style = {
  position: 'absolute' as 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: 400,
  bgcolor: 'background.paper',
  boxShadow: 24,
  p: 4,
  borderRadius: 2,
};

const PhoneInputModal = ({ open, onClose }: PhoneInputModalProps) => {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');

  const handleSendVerification = async () => {
    setError('');
    try {
      localStorage.setItem('adminEmail', email);
      const response = await axios.post(
        'http://localhost:8080/api/send-verification',
        { email },
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );
      if (response.status === 200) {
        navigate('/verify');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to send verification code');
    }
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      aria-labelledby="phone-modal-title"
      aria-describedby="phone-modal-description"
    >
      <Box sx={style}>
        <Typography id="phone-modal-title" variant="h6" component="h2" gutterBottom>
          Enter Admin Email
        </Typography>
        {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
        <TextField
          fullWidth
          label="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          error={!!error}
          helperText={error}
          sx={{ mb: 2 }}
          placeholder="Enter your email"
          type="email"
        />
        <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
          <Button onClick={onClose} variant="outlined">
            Cancel
          </Button>
          <Button
            onClick={handleSendVerification}
            variant="contained"
          >
            Send Verification
          </Button>
        </Box>
      </Box>
    </Modal>
  );
};

export default PhoneInputModal;
