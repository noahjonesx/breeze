import React, { useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  ActivityIndicator,
  ScrollView,
  Pressable,
} from "react-native";
import * as Clipboard from "expo-clipboard";
import { usePushNotifications } from "../hooks/usePushNotifications";

export default function HomeScreen() {
  const { expoPushToken, error } = usePushNotifications();
  const [copied, setCopied] = useState(false);

  async function copyToken() {
    if (!expoPushToken) return;
    await Clipboard.setStringAsync(expoPushToken);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <ScrollView contentContainerStyle={styles.container}>
      <Text style={styles.title}>Breeze</Text>
      <Text style={styles.subtitle}>Gazebo wind alerts</Text>

      <View style={styles.card}>
        <Text style={styles.cardLabel}>Your push token</Text>
        {!expoPushToken && !error && <ActivityIndicator style={styles.spinner} />}
        {error && <Text style={styles.error}>{error}</Text>}
        {expoPushToken && (
          <>
            <Text style={styles.token} numberOfLines={2}>
              {expoPushToken}
            </Text>
            <Pressable style={styles.button} onPress={copyToken}>
              <Text style={styles.buttonText}>{copied ? "Copied!" : "Copy token"}</Text>
            </Pressable>
          </>
        )}
      </View>

      <View style={styles.card}>
        <Text style={styles.cardLabel}>How it works</Text>
        <Text style={styles.info}>
          Your server polls wind forecasts every 15 minutes. You'll receive a
          notification ~2 hours before wind reaches your threshold, and again at
          ~30 minutes out.
        </Text>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flexGrow: 1,
    backgroundColor: "#0f172a",
    alignItems: "center",
    paddingTop: 80,
    paddingHorizontal: 24,
    paddingBottom: 40,
  },
  title: {
    fontSize: 48,
    fontWeight: "700",
    color: "#e2e8f0",
    letterSpacing: -1,
  },
  subtitle: {
    fontSize: 16,
    color: "#64748b",
    marginTop: 4,
    marginBottom: 40,
  },
  card: {
    width: "100%",
    backgroundColor: "#1e293b",
    borderRadius: 16,
    padding: 20,
    marginBottom: 16,
  },
  cardLabel: {
    fontSize: 12,
    fontWeight: "600",
    color: "#64748b",
    textTransform: "uppercase",
    letterSpacing: 1,
    marginBottom: 12,
  },
  spinner: {
    marginVertical: 8,
  },
  token: {
    fontSize: 12,
    color: "#94a3b8",
    fontFamily: "monospace",
    marginBottom: 12,
  },
  button: {
    backgroundColor: "#3b82f6",
    borderRadius: 10,
    paddingVertical: 10,
    alignItems: "center",
  },
  buttonText: {
    color: "#fff",
    fontWeight: "600",
    fontSize: 14,
  },
  info: {
    fontSize: 14,
    color: "#94a3b8",
    lineHeight: 22,
  },
  error: {
    fontSize: 13,
    color: "#f87171",
  },
});
