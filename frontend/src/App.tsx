import React from 'react';
import { Container, Box, Typography, CssBaseline } from '@mui/material';
import ContactForm from './components/ContactForm';
import { ThemeProvider, createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
    background: {
      default: '#98999e'
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ flexGrow: 1 }}>
        <Container maxWidth="lg">
          <Box sx={{ my: 4, textAlign: 'center' }}>
            <img src="/assets/chanterelle-logo.png" alt="Chanterelle Band Logo" style={{ maxWidth: '300px', height: 'auto', marginBottom: '1rem' }} />
            <Typography variant="h4" component="h1" gutterBottom align="center">
              A band based in Central Vermont
            </Typography>
            <Typography variant="h6" component="h2" paragraph align="center">
              We're excited to hear from you!
            </Typography>
            <ContactForm />
          </Box>
        </Container>
      </Box>
    </ThemeProvider>
  );
}

export default App;
