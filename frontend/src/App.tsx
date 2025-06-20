import React from 'react';
import { Container, Box, Typography, CssBaseline, Link } from '@mui/material';
import ContactForm from './components/ContactForm';
import { ThemeProvider, createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2'
    },
    secondary: {
      main: '#dc004e',
    },
    background: {
      default: '#fff', //'#98999e'
      paper: '#98999e'
    },
  },
});

function App() {
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
            <Box sx={{ m: 3, borderRadius: '40px', background: theme.palette.background.paper }}>
            <img src="/assets/chanterelle-logo.png" alt="Chanterelle Band Logo" style={{ maxWidth: '300px', height: 'auto', marginBottom: '1rem' }} />
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
            <Box sx={{ maxWidth: '500px', width: '100%', mt: 4 }}>
              <ContactForm />
            </Box>
          </Box>
        </Container>
      </Box>
    </ThemeProvider>
  );
}

export default App;
