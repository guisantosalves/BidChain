package com.bidchain.app

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableIntStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import com.bidchain.app.ui.theme.BidChainTheme
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue

/*
   Scaffold é o layout base do Material 3. Ele gerencia topbar, bottombar, FAB, etc. O innerPadding é o
   padding que ele calcula pra teu conteúdo não ficar atrás da status bar ou navigation bar (por causa
   do enableEdgeToEdge()).
 */
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            BidChainTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    CounterScreen(
                        modifier = Modifier.padding(innerPadding)
                    )
                }
            }
        }
    }
}

@Composable
fun CounterScreen(modifier: Modifier){
    var count by remember { mutableIntStateOf(0) }

    Column(modifier = modifier.fillMaxSize(), horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center){
        Text(text = "Count: $count")
        Button(onClick = {count++}) {
            Text("Increment")
        }
    }
}