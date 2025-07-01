import React from 'react';
import { Button, Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import PhoneInputModal from './PhoneInputModal';

const LoginButton = () => {
  const navigate = useNavigate();
  const [openModal, setOpenModal] = React.useState(false);

  const handleLogin = () => {
    setOpenModal(true);
  };

  return (
    <>
      <Button
        variant="contained"
        color="primary"
        onClick={handleLogin}
        sx={{ mt: 2 }}
      >
        Admin Login
      </Button>
      <PhoneInputModal open={openModal} onClose={() => setOpenModal(false)} />
    </>
  );
};

export default LoginButton;
