import React from 'react';
import { Button } from '@mui/material';
import LockIcon from '@mui/icons-material/Lock';
import EmailInputModal from './EmailInputModal';

const LoginButton = () => {
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
        <LockIcon />
      </Button>
      <EmailInputModal open={openModal} onClose={() => setOpenModal(false)} />
    </>
  );
};

export default LoginButton;
