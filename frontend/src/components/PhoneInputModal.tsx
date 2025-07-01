import React, { useState } from 'react';
import { Modal, Box, TextField, Button, Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

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

interface PhoneInputModalProps {
  open: boolean;
  onClose: () => void;
}

const PhoneInputModal = ({ open, onClose }: PhoneInputModalProps) => {
  const navigate = useNavigate();
  const [phoneNumber, setPhoneNumber] = useState('');
  const [error, setError] = useState('');

  const handleSendVerification = async () => {
    setError('');
    try {
      localStorage.setItem('adminPhoneNumber', phoneNumber);
      const response = await axios.post('http://localhost:8080/api/send-verification', { phoneNumber });
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
          Enter Phone Number
        </Typography>
        <TextField
          fullWidth
          label="Phone Number"
          value={phoneNumber}
          onChange={(e) => setPhoneNumber(e.target.value)}
          error={!!error}
          helperText={error}
          sx={{ mb: 2 }}
          placeholder="e.g., +18025551234"
        />
        <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
          <Button onClick={onClose} variant="outlined">
            Cancel
          </Button>
          <Button
            onClick={handleSendVerification}
            variant="contained"
            disabled={!phoneNumber.startsWith('+1') || phoneNumber.length !== 12}
          >
            Send Verification
          </Button>
        </Box>
      </Box>
    </Modal>
  );
};

export default PhoneInputModal;
