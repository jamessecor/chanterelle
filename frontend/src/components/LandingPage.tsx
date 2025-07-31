import React, { useEffect } from 'react';
import { Container, Box, Typography, CssBaseline, Link } from '@mui/material';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import ContactForm from './ContactForm';
import imgUrl from '../../assets/chanterelle-logo.png';
import LoginButton from './LoginButton';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2'
    },
    secondary: {
      main: '#dc004e',
    },
    background: {
      default: '#fff',
      paper: '#98999e'
    },
  },
});

const LandingPage = () => {
  useEffect(() => {
    const logoObject = document.getElementById('logoImage');
    logoObject?.setAttribute('src', imgUrl);
  }, [imgUrl]);

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh',
        width: '100%',
        px: 2
      }}>
        <Container maxWidth="sm" sx={{ textAlign: 'center', width: '100%' }}>
          <Box sx={{ textAlign: 'center', width: '100%' }}>
            <Box sx={{ p: { xs: 0, md: 3 }, my: 3, borderRadius: { xs: '10px', md: '40px' }, background: theme.palette.background.paper }}>
              <img id="logoImage" alt="Chanterelle Band Logo" style={{ maxWidth: '300px', height: 'auto', marginBottom: '1rem' }} />
            </Box>
            <Typography variant="h4" component="h1" gutterBottom align="center">
              A band based in Central Vermont
            </Typography>
            <Typography variant="h6" component="h2" paragraph align="center">
              We're excited to hear from you!
            </Typography>
            <Typography variant="body1" align="center" sx={{ mb: 2 }}>
              Check out our music on{' '}
              <Link href="https://chanterellevt.bandcamp.com" target="_blank" rel="noopener noreferrer" color="primary">
                Bandcamp
              </Link>
            </Typography>
            <Box sx={{ 
              maxWidth: '500px', 
              width: '100%', 
              mt: 4,
              mx: 'auto'  
            }}>
              <ContactForm />
            </Box>
            <LoginButton />
          </Box>
        </Container>
      </Box>
    </ThemeProvider>
  );
};

export default LandingPage;
